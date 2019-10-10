package main

// The SkyLB management UI prototype.
// Functions:
// Set service as offline ("lame duck") or online.
// List services.
//
// Http request format:
// baseurl/svc?action=online&eps=ip1:port1,ip2:port2,...
// baseurl/svc?action=offline&eps=ip1:port1,ip2:port2,...

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	prom "github.com/prometheus/client_golang/prometheus"

	"binchencoder.com/letsgo"
	"binchencoder.com/letsgo/metrics"
	js "binchencoder.com/letsgo/strings"
	"binchencoder.com/skylb-api/lameduck"
	"binchencoder.com/skylb-api/prefix"
	"binchencoder.com/skylb/cmd/grpchealth"
	"binchencoder.com/skylb/cmd/webserver/svclist"
	"binchencoder.com/skylb/hub"
)

const (
	ACTION = "action"

	// Short for "endpoints"
	EPS = "eps"
	// Short for "service"
	SVC = "svc"
)

var (
	hostPort          = flag.String("host-port", ":8090", "The web server host:port")
	scrapePort        = flag.Int("scrape-port", 8092, "The port to listen on for HTTP requests.")
	checkSvcInterval  = flag.Duration("check-svc-interval", 10*time.Minute, "How often to check service health")
	enableHealthCheck = flag.Bool("enable_health_check", false, "Enable health check or not.")

	etcdCli etcd.KeysAPI

	svcHealthTimestamp = prom.NewGauge(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "web",
			Name:      "service_health_timestamp",
			Help:      "The timestamp of last successful service health check.",
		},
	)
)

func init() {
	prom.MustRegister(svcHealthTimestamp)
}

func main() {
	letsgo.Init()

	etcdCli = hub.CreateEtcdClient(*hub.EtcdEndpoints, true)
	prefix.Init(etcdCli)

	http.HandleFunc("/svc", handleSvcRequest)
	http.HandleFunc("/api/v1/svc/list", handleSvcList)

	go metrics.StartPrometheusServer(*scrapePort)

	if *enableHealthCheck {
		go checkServices(etcdCli)
	}

	glog.Infof("Starting web server at %s\n", *hostPort)
	err := http.ListenAndServe(*hostPort, nil)
	if err != nil {
		glog.Fatal("ListenAndServe: ", err)
	}
}

func handleSvcList(w http.ResponseWriter, r *http.Request) {
	listSvc(w)
}

func handleSvcRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	act := query.Get(ACTION)
	svc := query.Get(SVC)
	if svc == "" {
		// TODO(zhwang): disable dymmy-service. At present, the service name
		//               is not in the REST API.
		//               To make the tool work we temporarilly use the dummy
		//               service name.
		svc = "dummy-service"
	}
	eps := query.Get(EPS)
	var err error
	switch strings.ToLower(act) {
	case "offline":
		err = offline(svc, eps)
	case "online":
		err = online(svc, eps)
	case "list":
		listSvc(w)
		return
	case "listdep":
		listDep(w)
		return
	default:
		respond(w, fmt.Sprintf("Unsupported action %s", act))
		return
	}
	if nil != err {
		respond(w, fmt.Sprintf("%v", err))
	} else {
		respond(w, "OK")
	}
}

// listDep lists services and their users (also services).
func listDep(w http.ResponseWriter) {
	// TODO(fuyc): implement it.
}

// listSvc simply lists all services
func listSvc(w http.ResponseWriter) {
	list, err := svclist.ListServices(etcdCli)
	if nil != err {
		respond(w, fmt.Sprintf("%v", err))
		return
	}
	// Pretty-print as json.
	b, err := json.MarshalIndent(&list, "", "\t")
	if nil != err {
		respond(w, fmt.Sprintf("%v", err))
		return
	}
	w.Write(b)
}

func respond(w http.ResponseWriter, msg string) {
	if _, err := io.WriteString(w, msg); err != nil {
		glog.Warningf("Failed to write response, %#v", err)
	}
}

func parseEndpoints(epsStr string) ([]hub.ServiceEndpoint, error) {
	eps := js.CsvToSlice(epsStr)
	if len(eps) == 0 {
		return nil, fmt.Errorf("No valid endpoints: %s", epsStr)
	}
	addrs := make([]hub.ServiceEndpoint, 0, 1000)
	for _, ep := range eps {
		ep = strings.TrimSpace(ep)
		host, port, err := net.SplitHostPort(ep)
		if err != nil {
			glog.Warning(err)
			continue
		}
		portnum, err := net.LookupPort("tcp", port)
		if err != nil {
			glog.Warning(err)
			continue
		}
		addr := hub.ServiceEndpoint{
			IP:   host,
			Port: int32(portnum),
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		return nil, errors.New("No valid endpoints")
	}
	return addrs, nil
}

func offline(svc, epsStr string) error {
	addrs, err := parseEndpoints(epsStr)
	if nil != err {
		return err
	}
	for _, addr := range addrs {
		if err := lameduck.SetLameDuckMode(etcdCli, svc, lameduck.HostPort(addr.IP, fmt.Sprintf("%d", addr.Port))); nil != err {
			glog.Warningf("Set lame duck err:%#v", err)
		}
	}
	return nil
}

func online(svc, epsStr string) error {
	addrs, err := parseEndpoints(epsStr)
	if nil != err {
		return err
	}
	for _, addr := range addrs {
		if err := lameduck.UnsetLameDuckMode(etcdCli, svc, lameduck.HostPort(addr.IP, fmt.Sprintf("%d", addr.Port))); nil != err {
			// Unsetting an inexistent lame duck should not
			// need to warn, just info is enough.
			glog.V(2).Infof("Unset lame duck err:%#v", err)
		}
	}
	return nil
}

func checkServices(etcdCli etcd.KeysAPI) {
	time.Sleep(3 * time.Second)
	for {
		grpchealth.CheckAllServices(etcdCli)
		svcHealthTimestamp.Set(float64(time.Now().Unix()))
		time.Sleep(*checkSvcInterval)
	}
}
