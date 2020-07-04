package grpchealth

import (
	"flag"
	"fmt"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	prom "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	hpb "google.golang.org/grpc/health/grpc_health_v1"

	vexpb "github.com/binchencoder/gateway-proto/data"
	"github.com/binchencoder/letsgo/service/naming"
	skylbclient "github.com/binchencoder/skylb-api/client"
	jh "github.com/binchencoder/skylb-api/internal/health" // TODO(fuyc): may refactor the package.
	skypb "github.com/binchencoder/skylb-api/proto"
	"github.com/binchencoder/skylb/cmd/webserver/svclist"
)

const (
	// Various ways to check service:
	// Renew skylb client.
	renew = "renew"
	// Reuse skylb client.
	reuse = "reuse"
	// Raw grpc.
	rawGrpc = "grpc"
)

var (
	healthCheckTimeout = jh.HealthCheckTimeout
	enableRenewSkylb   = flag.Bool("enable-renew-skylb", true, "Whether check services using newly created skylb client")
	enableRawGrpc      = flag.Bool("enable-raw-grpc", true, "Whether check services via raw grpc, instead of skylb")

	// Map of service id to health client.
	sidClients = make(map[vexpb.ServiceId]hpb.HealthClient, 50)

	svcHealthGauge = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "web",
			Name:      "service_health_gauge",
			Help:      "Service health status.",
		},
		[]string{"service", "grpc_code", "type"},
	)

	svcHealthDialCounts = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "web",
			Name:      "service_health_dial_counts",
			Help:      "Service health dial counts.",
		},
		[]string{"service", "type"},
	)
	svcHealthStartCounts = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "web",
			Name:      "service_health_start_counts",
			Help:      "Service health start counts.",
		},
		[]string{"service", "type"},
	)

	svcHealthSuccessRate = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "web",
			Name:      "service_health_success_rate",
			Help:      "grpc health check success rate.",
		},
		[]string{"service", "type"},
	)

	svcHealthDialLatencyHistogram = prom.NewHistogramVec(
		prom.HistogramOpts{
			Namespace: "skylb",
			Subsystem: "web",
			Name:      "service_health_dial_latency",
			Help:      "Histogram of dial latency of grpc health checking.",
			Buckets:   prom.DefBuckets,
		},
		[]string{"service", "type"},
	)
	svcHealthLatencyHistogram = prom.NewHistogramVec(
		prom.HistogramOpts{
			Namespace: "skylb",
			Subsystem: "web",
			Name:      "health_check_latency",
			Help:      "grpc health check latency.",
			Buckets:   prom.DefBuckets,
		},
		[]string{"service", "type"},
	)
)

func init() {
	prom.MustRegister(svcHealthGauge)
	prom.MustRegister(svcHealthDialCounts)
	prom.MustRegister(svcHealthStartCounts)
	prom.MustRegister(svcHealthSuccessRate)
	prom.MustRegister(svcHealthDialLatencyHistogram)
	prom.MustRegister(svcHealthLatencyHistogram)
}

// CheckAllServices checks all grpc services against their grpc health interface.
func CheckAllServices(etcdCli etcd.KeysAPI) {
	sepsArr, err := svclist.ListServices(etcdCli)
	if nil != err {
		glog.Error(err)
		return
	}
	svcHealthGauge.Reset()
	for _, seps := range sepsArr {
		if 0 == len(seps.InstEndpoints) {
			glog.V(3).Infof("No instances for service: %s", seps.Spec.ServiceName)
			svcHealthGauge.WithLabelValues(seps.Spec.ServiceName, "NA", "NA").Set(0)
			svcHealthSuccessRate.WithLabelValues(seps.Spec.ServiceName, "NA").Set(0)
			continue
		}

		// TODO(fuyc): run health checks concurrently for different services,
		//             or run with a bounded fan out.
		svcName := seps.Spec.ServiceName
		serviceId, err := naming.ServiceNameToId(svcName)
		if nil != err {
			glog.Error(err)
			continue
		}

		if *enableRawGrpc {
			CheckServiceRawGrpc(serviceId, svcName, seps)
		}
		//TODO(fuyc): compare with ep count
		CheckServiceReuse(serviceId, svcName, len(seps.InstEndpoints))
		if *enableRenewSkylb {
			CheckService(serviceId, svcName, len(seps.InstEndpoints))
		}
	}
}

// CheckService checks grpc services by reusing skylb client.
func CheckServiceReuse(serviceId vexpb.ServiceId, serviceName string, count int) {

	healthCli := sidClients[serviceId]
	if nil == healthCli {
		// Initialize gRPC service client.
		skycli := skylbclient.NewServiceCli(vexpb.ServiceId_SKYLB_FRONTEND)

		skycli.EnableHistogram()

		skycli.EnableFailFast()

		// Resolve service
		demoSpec := skylbclient.NewServiceSpec("", serviceId, "")
		skycli.Resolve(demoSpec)

		svcHealthDialCounts.WithLabelValues(serviceName, reuse).Inc()
		start := time.Now()
		skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
			switch spec.String() {
			case demoSpec.String():
				if nil != conn {
					healthCli = hpb.NewHealthClient(conn)
				}
			}
		})
		svcHealthDialLatencyHistogram.WithLabelValues(serviceName, reuse).Observe(float64(time.Since(start).Seconds()))
		sidClients[serviceId] = healthCli
		// Keep the connection without shutting down.
	}

	checkServiceSkylb(healthCli, count, serviceName, reuse)
}

// CheckService checks grpc services using newly created skylb client.
func CheckService(serviceId vexpb.ServiceId, serviceName string, count int) {

	// Initialize gRPC service client.
	skycli := skylbclient.NewServiceCli(vexpb.ServiceId_SKYLB_FRONTEND)

	skycli.EnableHistogram()

	skycli.EnableFailFast()

	// Resolve service
	demoSpec := skylbclient.NewServiceSpec("", serviceId, "")
	skycli.Resolve(demoSpec)

	var healthCli hpb.HealthClient
	svcHealthDialCounts.WithLabelValues(serviceName, renew).Inc()
	start := time.Now()
	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			if nil != conn {
				healthCli = hpb.NewHealthClient(conn)
			}
		}
	})
	svcHealthDialLatencyHistogram.WithLabelValues(serviceName, renew).Observe(float64(time.Since(start).Seconds()))
	defer skycli.Shutdown()

	checkServiceSkylb(healthCli, count, serviceName, renew)
}

func checkServiceSkylb(healthCli hpb.HealthClient, count int, serviceName string, checkType string) {
	if nil == healthCli {
		glog.Infof("Resolved 0 instances of service %s", serviceName)
		svcHealthGauge.WithLabelValues(serviceName, "NA", checkType).Set(0)
		svcHealthSuccessRate.WithLabelValues(serviceName, checkType).Set(0)
		return
	}

	glog.V(3).Infof("===== [%s] Checking %s =====", checkType, serviceName)
	success := 0
	for i := 0; i < count; i++ {
		req := hpb.HealthCheckRequest{}
		ctx, cancel := context.WithTimeout(context.Background(), *healthCheckTimeout)
		start := time.Now()
		svcHealthStartCounts.WithLabelValues(serviceName, checkType).Inc()
		resp, err := healthCli.Check(ctx, &req, grpc.FailFast(false))
		svcHealthLatencyHistogram.WithLabelValues(serviceName, checkType).Observe(time.Since(start).Seconds())
		glog.V(3).Infof("Reply: %#v", resp)
		if err != nil {
			cancel()
			svcHealthGauge.WithLabelValues(serviceName, grpc.Code(err).String(), checkType).Inc()
			if !jh.IsSafeError(err) {
				glog.Errorf("Failed to check %s, %v", serviceName, err)
				continue
			}
		} else {
			svcHealthGauge.WithLabelValues(serviceName, "OK", checkType).Inc()
		}
		success++
	}
	svcHealthSuccessRate.WithLabelValues(serviceName, checkType).Set(float64(success) / float64(count))
	glog.V(3).Infof("===== Finished checking %s =====", serviceName)
}

// CheckServiceRawGrpc check services using raw grpc.
func CheckServiceRawGrpc(serviceId vexpb.ServiceId, serviceName string, seps *skypb.ServiceEndpoints) {
	glog.V(3).Infof("===== [grpc] Checking %s =====", serviceName)
	success := 0
	size := len(seps.InstEndpoints)
	for _, iep := range seps.InstEndpoints {
		func() {
			start := time.Now()
			addr := fmt.Sprintf("%s:%d", iep.Host, iep.Port)
			svcHealthDialCounts.WithLabelValues(serviceName, rawGrpc).Inc()
			conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(*healthCheckTimeout))
			svcHealthDialLatencyHistogram.WithLabelValues(serviceName, rawGrpc).Observe(float64(time.Since(start).Seconds()))
			if nil != err {
				glog.Errorf("Error dialling service %s at %s", serviceName, addr)
				return
			}
			defer conn.Close()

			healthCli := hpb.NewHealthClient(conn)
			req := hpb.HealthCheckRequest{}
			ctx, cancel := context.WithTimeout(context.Background(), *healthCheckTimeout)
			start = time.Now()
			svcHealthStartCounts.WithLabelValues(serviceName, rawGrpc).Inc()
			resp, err := healthCli.Check(ctx, &req, grpc.FailFast(false))
			svcHealthLatencyHistogram.WithLabelValues(serviceName, rawGrpc).Observe(time.Since(start).Seconds())
			glog.V(3).Infof("Reply: %#v", resp)
			if nil != err {
				cancel()
				svcHealthGauge.WithLabelValues(serviceName, grpc.Code(err).String(), rawGrpc).Inc()
				if !jh.IsSafeError(err) {
					glog.Errorf("Failed to check %s, %v", serviceName, err)
					return
				}
			} else {
				svcHealthGauge.WithLabelValues(serviceName, "OK", rawGrpc).Inc()
			}
			success++
		}()
	}
	svcHealthSuccessRate.WithLabelValues(serviceName, rawGrpc).Set(float64(success) / float64(size))
	glog.V(3).Infof("===== Finished checking %s =====", serviceName)
}
