package rpc

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/golang/glog"
	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/peer"

	"binchencoder.com/skylb-api/lameduck"
	pb "binchencoder.com/skylb-api/proto"
	"binchencoder.com/skylb/hub"
)

var (
	flagNotifyTimeout      = flag.Duration("endpoints-notify-timeout", 10*time.Second, "The timeout to notify client endpoints update")
	flagAutoDisconnTimeout = flag.Duration("auto-disconn-timeout", 5*time.Minute, "The timeout to automatically disconnect the resolve RPC")

	activeObserverGauge = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "active_observer_gauge",
			Help:      "SkyLB active observer gauge.",
		},
		[]string{"service"},
	)
	activeReporterGauge = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "active_reporter_gauge",
			Help:      "SkyLB active reporter gauge.",
		},
		[]string{"endpoint"},
	)
	addObserverFailCounts = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "add_observer_fail_counts",
			Help:      "SkyLB observer rpc counts.",
		},
		[]string{"service"},
	)
	observeRpcCounts = prom.NewCounter(
		prom.CounterOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "observe_rpc_counts",
			Help:      "SkyLB observer rpc counts.",
		},
	)
	reportLoadRpcCounts = prom.NewCounter(
		prom.CounterOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "report_load_rpc_counts",
			Help:      "SkyLB report load rpc counts.",
		},
	)
	reportLoadCounts = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "report_load_counts",
			Help:      "SkyLB report load counts.",
		},
		[]string{"service"},
	)
	initReportLoadCounts = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "init_report_load_counts",
			Help:      "SkyLB init report load counts.",
		},
		[]string{"service"},
	)
	// To test this metric, enable flag:
	// --auto-disconn-timeout=2s
	autoDisconnCounts = prom.NewCounter(
		prom.CounterOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "auto_disconn_counts",
			Help:      "SkyLB auto disconnect counts.",
		},
	)
	// To test this metric, enable flag:
	// --endpoints-notify-timeout=1ns
	notifyTimeoutCounts = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "notify_timeout_counts",
			Help:      "Notify client endpoints update timeout counts.",
		},
		[]string{"caller_service", "caller_addr"},
	)
	notifyChanUsageHistogram = prom.NewHistogram(
		prom.HistogramOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "notify_chan_usage",
			Help:      "The usage rate of the notify channel.",
			Buckets:   prom.LinearBuckets(0, 0.1, 10),
		},
	)
)

func init() {
	prom.MustRegister(activeObserverGauge)
	prom.MustRegister(activeReporterGauge)
	prom.MustRegister(addObserverFailCounts)
	prom.MustRegister(autoDisconnCounts)
	prom.MustRegister(initReportLoadCounts)
	prom.MustRegister(observeRpcCounts)
	prom.MustRegister(reportLoadCounts)
	prom.MustRegister(reportLoadRpcCounts)
	prom.MustRegister(notifyChanUsageHistogram)
	prom.MustRegister(notifyTimeoutCounts)
}

// Struct skylbServer implements interface pb.SkylbServer.
type skylbServer struct {
	epsHub hub.EndpointsHub
}

func (ss *skylbServer) Resolve(req *pb.ResolveRequest, stream pb.Skylb_ResolveServer) error {
	glog.V(4).Infof("Caller service %d\n", req.CallerServiceId)
	observeRpcCounts.Inc()

	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return errors.New("Failed to get peer client info from context.")
	}

	if len(req.Services) == 0 {
		return errors.New("No service spec found.")
	}

	for _, svc := range req.Services {
		ss.epsHub.TrackServiceGraph(req, svc, p.Addr)
	}

	defer func() {
		for _, svc := range req.Services {
			ss.epsHub.UntrackServiceGraph(req, svc, p.Addr)
		}
	}()

	notiCh, err := ss.epsHub.AddObserver(req.Services, p.Addr.String(), req.ResolveFullEndpoints)
	if err != nil {
		for _, s := range req.Services {
			label := fmt.Sprintf("%s.%s", s.Namespace, s.ServiceName)
			addObserverFailCounts.WithLabelValues(label).Inc()
		}

		glog.Infof("Failed to register caller service ID %d client %s to observe services, %+v", req.CallerServiceId, p.Addr.String(), err)
		return err
	}
	defer func() {
		glog.Infof("Stop observing services for caller service ID %d client %s", req.CallerServiceId, p.Addr.String())
		ss.epsHub.RemoveObserver(req.Services, p.Addr.String())
		for _, spec := range req.Services {
			label := fmt.Sprintf("%s.%s", spec.Namespace, spec.ServiceName)
			activeObserverGauge.WithLabelValues(label).Dec()
		}
	}()

	maxIds := map[string]int64{}
	for _, spec := range req.Services {
		maxIds[spec.String()] = 0

		label := fmt.Sprintf("%s.%s", spec.Namespace, spec.ServiceName)
		activeObserverGauge.WithLabelValues(label).Inc()
		glog.Infof("Registered caller service ID %d client %s to observe service %s.%s on port name %q", req.CallerServiceId, p.Addr.String(), spec.Namespace, spec.ServiceName, spec.PortName)
	}

	timer := time.NewTimer(*flagAutoDisconnTimeout + time.Duration(rand.Int63n(int64(*flagAutoDisconnTimeout))))
	for {
		select {
		case <-timer.C:
			glog.Infoln("Auto disconnect with client")
			autoDisconnCounts.Inc()
			return nil
		case updates, ok := <-notiCh:
			if !ok {
				// Channel has been closed.
				return nil
			}

			notifyChanUsageHistogram.Observe(float64(len(notiCh)) / hub.ChanCapMultiplication / float64(len(req.Services)))

			maxId := maxIds[updates.Endpoints.Spec.String()]
			if updates.Id < maxId {
				// Skip the old updates.
				continue
			} else {
				maxIds[updates.Endpoints.Spec.String()] = updates.Id
			}

			eps := updates.Endpoints

			if glog.V(3) {
				var buf bytes.Buffer
				for i, iep := range eps.InstEndpoints {
					if i > 0 {
						(&buf).WriteString(", ")
					}
					(&buf).WriteString(fmt.Sprintf("[%s]%s:%d", opToString(iep.Op), iep.Host, iep.Port))
				}
				if req.ResolveFullEndpoints {
					glog.Infof("Full endpoints of service %s for caller service ID %d client %s: %s.",
						eps.Spec.ServiceName, req.CallerServiceId, p.Addr.String(), buf.String())
				} else {
					glog.Infof("Endpoints changed for caller service ID %d client %s with updates %s.", req.CallerServiceId, p.Addr.String(), buf.String())
				}
			} else {
				if req.ResolveFullEndpoints {
					glog.Infof("Send full endpoints of service %s for caller service ID %d client %s.",
						eps.Spec.ServiceName, req.CallerServiceId, p.Addr.String())
				} else {
					glog.Infof("Endpoints changed for caller service ID %d client %s.", req.CallerServiceId, p.Addr.String())
				}
			}

			resp := pb.ResolveResponse{
				SvcEndpoints: &pb.ServiceEndpoints{
					Spec:          eps.Spec,
					InstEndpoints: eps.InstEndpoints,
				},
			}

			errCh := make(chan error, 1)
			go func(ch chan<- error) {
				if err := stream.Send(&resp); err != nil {
					ch <- err
					return
				}
				ch <- nil
			}(errCh)

			t := time.NewTimer(*flagNotifyTimeout)
			select {
			case <-t.C:
				// Discard the current gRPC stream if timeout.
				glog.Errorf("Time out to send endpoints update to caller service ID %d client %s, abandon the stream.", req.CallerServiceId, p.Addr.String())
				// It's OK to record p.Addr.String in label value, since such events should be rare, and will not accumulate too much data.
				notifyTimeoutCounts.WithLabelValues(fmt.Sprintf("%d", req.CallerServiceId), p.Addr.String()).Inc()
				return errors.New("time out to send endpoints update to client")
			case err := <-errCh:
				if !t.Stop() {
					<-t.C
				}
				if err != nil {
					glog.Errorf("Failed to send endpoints update to caller service ID %d client %s, abandon the stream, %+v.", req.CallerServiceId, p.Addr.String(), err)
					return err
				}
			}
		}
	}
}

func (ss *skylbServer) ReportLoad(stream pb.Skylb_ReportLoadServer) error {
	reportLoadRpcCounts.Inc()

	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return errors.New("failed to get peer info from context")
	}

	// Extract the peer's host name.
	host := p.Addr.String()
	pos := strings.Index(host, ":")
	if pos > -1 {
		host = host[:pos]
	}

	glog.Infof("Start accepting load report from %s.", host)
	activeReporterGauge.WithLabelValues(host).Inc()
	defer func() {
		activeReporterGauge.WithLabelValues(host).Dec()
	}()

	first := true
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		label := fmt.Sprintf("%s.%s", req.Spec.Namespace, req.Spec.ServiceName)
		reportLoadCounts.WithLabelValues(label).Inc()

		// Replace host name if fixed_host has been specified.
		h := host
		if req.FixedHost != "" {
			h = req.FixedHost
			glog.V(4).Infof("Use fixed host %s instead of %s", h, host)
		}

		if first {
			glog.V(3).Infof("Received init load report from %s at %s", label, p.Addr.String())
			initReportLoadCounts.WithLabelValues(label).Inc()

			// When the service with weights is turned off, the service
			// is restarted in less than 10 seconds, especially if
			// the weights are modified. If only the epsHub.UpsertEndpoint
			// method is used, the weight level is not modified.
			// purely Just to prevent this issue.
			if err := ss.epsHub.InsertEndpoint(req.Spec, h, req.Port, req.Weight); err != nil {
				glog.Errorf("Failed to update etcd entry for endpoint %s:%d, closing the report stream.", h, req.Port)
				return err
			}

			first = false
		}

		// Block the heart beat if the server is lame duck.
		ep := lameduck.HostPort(h, fmt.Sprintf("%d", req.Port))
		if lameduck.IsLameduckMode(ep) {
			glog.V(4).Infof("Received load report from %s:%d, masked", h, req.Port)
			continue
		}
		glog.V(4).Infof("Received load report from %s:%d.", h, req.Port)

		if err := ss.epsHub.UpsertEndpoint(req.Spec, h, req.Port, req.Weight); err != nil {
			glog.Errorf("Failed to update etcd entry for endpoint %s:%d, closing the report stream.", h, req.Port)
			return err
		}
	}
}

func (ss *skylbServer) AttachForDiagnosis(stream pb.Skylb_AttachForDiagnosisServer) error {
	// TODO(fuyc): implement it.
	return nil
}

// NewSkylbServer creates and returns a new SkyLB gRPC server.
func NewSkylbServer() pb.SkylbServer {
	return &skylbServer{
		epsHub: hub.Init(),
	}
}

func opToString(op pb.Operation) string {
	switch op {
	case pb.Operation_Add:
		return "ADD"
	case pb.Operation_Delete:
		return "DELETE"
	}
	return ""
}
