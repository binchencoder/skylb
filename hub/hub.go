package hub

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	api "k8s.io/api/core/v1"

	"github.com/binchencoder/letsgo/strings"
	jsync "github.com/binchencoder/letsgo/sync"
	"github.com/binchencoder/skylb-api/lameduck"
	"github.com/binchencoder/skylb-api/prefix"
	pb "github.com/binchencoder/skylb-api/proto"
	"github.com/binchencoder/skylb-api/util"
)

const (
	TimestampKey = "timestamp"
	AddrKey      = "addr"
)

var (
	EtcdEndpoints    = flag.String("etcd-endpoints", "", "Comma separated etcd endpoints")
	graphKeyTTL      = flag.Duration("graph-key-ttl", 24*time.Hour, "The service graph key TTL")
	graphKeyInterval = flag.Duration("graph-key-interval", 2*time.Hour, "The service graph key update interval")
	withinK8s        = flag.Bool("within-k8s", false, "Whether SkyLB is running in kubernetes")

	getOpts      etcd.GetOptions
	setGraphOpts *etcd.SetOptions
	watchOpts    etcd.WatcherOptions

	nextUpdateId int64

	hub  *endpointsHub
	once sync.Once
)

func init() {
	getOpts = etcd.GetOptions{
		Recursive: true,
	}
	watchOpts = etcd.WatcherOptions{
		Recursive: true,
	}
	setGraphOpts = &etcd.SetOptions{
		TTL: *graphKeyTTL,
	}
}

type clientObject struct {
	spec        *pb.ServiceSpec
	clientAddr  string
	resolveFull bool
	notifyCh    chan<- *EndpointsUpdate
	stopCh      chan struct{}
}

// ServiceEndpoint represents a service endpoint.
// (A simplified version of pb.InstanceEndpoint)
type ServiceEndpoint struct {
	IP     string
	Port   int32
	Weight int32
}

func (se ServiceEndpoint) toString() string {
	return fmt.Sprintf("%s:%d", se.IP, se.Port)
}

type serviceEndpoints map[string]ServiceEndpoint

type serviceObject struct {
	jsync.RWLock

	spec      *pb.ServiceSpec
	endpoints serviceEndpoints
	observers []*clientObject
}

type serviceMap map[string]*serviceObject

// EndpointsHub defines the service endpoints hub based on etcd.
type EndpointsHub interface {
	// AddObserver adds an observer of the given service specs for the given
	// clientAddr. When service endpoints changed, it notifies the observer
	// through the returned channel.
	AddObserver(specs []*pb.ServiceSpec, clientAddr string, resolveFull bool) (<-chan *EndpointsUpdate, error)

	// RemoveObserver removes the observer for the given service specs for the
	// given clientAddr.
	RemoveObserver(specs []*pb.ServiceSpec, clientAddr string)

	// InsertEndpoint inserts a service with the given namespace and service name.
	InsertEndpoint(spec *pb.ServiceSpec, host string, port, weight int32) error

	// UpsertEndpoint inserts or update a service with the given namespace
	// and service name.
	UpsertEndpoint(spec *pb.ServiceSpec, host string, port, weight int32) error

	// TrackServiceGraph keeps track of dependency graph between clients and services.
	TrackServiceGraph(req *pb.ResolveRequest, callee *pb.ServiceSpec, callerAddr net.Addr)

	// UntrackServiceGraph stops tracking of dependency graph between clients and services.
	UntrackServiceGraph(req *pb.ResolveRequest, callee *pb.ServiceSpec, callerAddr net.Addr)
}

type endpointsHub struct {
	jsync.RWLock

	etcdKeyTtl time.Duration
	services   serviceMap
	etcdCli    etcd.KeysAPI

	graphKeys     map[string]struct{}
	graphKeysLock *sync.RWMutex
}

// InsertEndpoint inserts a service with the given namespace and service name.
func (eh *endpointsHub) InsertEndpoint(spec *pb.ServiceSpec, host string, port, weight int32) error {
	key := eh.calculateEndpointKey(spec.Namespace, spec.ServiceName, host, port)
	ctx := context.Background()
	return eh.setKey(ctx, key, spec, host, port, weight)
}

// UpsertEndpoint inserts or update a service with the given namespace
// and service name.
func (eh *endpointsHub) UpsertEndpoint(spec *pb.ServiceSpec, host string, port, weight int32) error {
	key := eh.calculateEndpointKey(spec.Namespace, spec.ServiceName, host, port)

	ctx := context.Background()
	err := eh.refreshKey(ctx, key)
	if err == nil {
		return nil
	}

	if e, ok := err.(etcd.Error); ok {
		switch e.Code {
		case etcd.ErrorCodeKeyNotFound:
			// Sometimes the key might be dropped or expired so that
			// refreshKey will fail.
			return eh.setKey(ctx, key, spec, host, port, weight)
		}
	}
	return err
}

func (eh *endpointsHub) fetchEndpoints(namespace, serviceName string) (*api.Endpoints, error) {
	endpoints := api.Endpoints{}

	key := eh.calculateKey(namespace, serviceName)
	resp, err := eh.etcdCli.Get(context.Background(), key, &getOpts)
	if err != nil {
		if e, ok := err.(etcd.Error); ok && e.Code == etcd.ErrorCodeKeyNotFound {
			glog.Warningf("Service %s.%s absent, return empty list.", namespace, serviceName)
			return &endpoints, nil
		}
		return nil, err
	}

	if resp.Node == nil {
		glog.Errorf("No endpoints found for service %s.%s", namespace, serviceName)
		return &endpoints, nil
	}

	if len(resp.Node.Nodes) == 0 {
		if resp.Node.Value != "" {
			if err := json.Unmarshal([]byte(resp.Node.Value), &endpoints); err != nil {
				return nil, err
			}
		}
		return &endpoints, nil
	}

	for i, node := range resp.Node.Nodes {
		if node.Value == "" {
			continue
		}
		eps := api.Endpoints{}
		if err := json.Unmarshal([]byte(node.Value), &eps); err != nil {
			return nil, err
		}
		if i == 0 {
			endpoints.Name = eps.Name
			endpoints.Namespace = eps.Namespace
			endpoints.ObjectMeta = eps.ObjectMeta
		}
		for labelKey, labelValue := range eps.Labels {
			if endpoints.Labels == nil {
				endpoints.Labels = make(map[string]string)
			}
			endpoints.Labels[labelKey] = labelValue
		}
		for _, subset := range eps.Subsets {
			endpoints.Subsets = append(endpoints.Subsets, subset)
		}
	}

	return &endpoints, nil
}

func (eh *endpointsHub) startMainWatcher() {
outerLoop:
	for {
		w := eh.etcdCli.Watcher(prefix.EndpointsKey, &watchOpts)

		// Watch etcd keys for all service endpoints and notify clients.
		for {
			resp, err := w.Next(context.Background())
			glog.V(4).Infof("Watched: %+v", resp)
			if err != nil {
				time.Sleep(time.Second)
				if e, ok := err.(etcd.Error); ok {
					if e.Code == etcd.ErrorCodeEventIndexCleared ||
						e.Code == etcd.ErrorCodeWatcherCleared {
						glog.Errorf("Abandon watcher, %v", err)
						continue outerLoop
					}
				}
				glog.Errorf("Failed to get next watch event, %v", err)
				continue
			}
			eh.extractUpdates(resp)
		}
	}
}

// startLameDuckWatcher starts a watcher to watch changes of lame duck.
func (eh *endpointsHub) startLameDuckWatcher() {
	// Load current lameduck endpoints.
	resp, err := eh.etcdCli.Get(context.Background(), prefix.LameduckKey, &etcd.GetOptions{Recursive: true})
	if err != nil {
		glog.Errorf("Failed to load lameduck instances with key prefix %s, %v", prefix.LameduckKey, err)
	} else {
		lameduck.ExtractLameduck(resp.Node)
	}

outerLoop:
	for {
		w := eh.etcdCli.Watcher(prefix.LameduckKey, &watchOpts)
		for {
			resp, err := w.Next(context.Background())
			glog.V(4).Infof("Watched lameduck change: %+v", resp)
			if err != nil {
				time.Sleep(time.Second)
				if e, ok := err.(etcd.Error); ok {
					if e.Code == etcd.ErrorCodeEventIndexCleared ||
						e.Code == etcd.ErrorCodeWatcherCleared {
						glog.Errorf("Abandon watcher, %v", err)
						continue outerLoop
					}
				}
				glog.Errorf("Failed to get next watch event, %v", err)
				continue
			}
			lameduck.ExtractLameduckChange(resp)
		}
	}
}

// SkyLB receives full set of current endpoints.
func (eh *endpointsHub) extractUpdates(resp *etcd.Response) {
	var key string
	switch resp.Action {
	case util.ActionCreate, util.ActionSet:
		key = path.Dir(resp.Node.Key)
	case util.ActionDelete, util.ActionExpire:
		key = path.Dir(resp.PrevNode.Key)
	default:
		glog.Errorf("Unexpected action %s, ignore.", resp.Action)
		return
	}

	eh.updateEndpoints(key)
}

func (eh *endpointsHub) updateEndpoints(key string) {
	var so *serviceObject
	eh.WithRLock(func() error {
		so, _ = eh.services[key]
		return nil
	})
	if so == nil {
		glog.V(3).Infof("serviceObject nil for key %#v", key)
		return
	}

	eps, err := eh.fetchEndpoints(so.spec.Namespace, so.spec.ServiceName)
	if err != nil {
		glog.Errorf("Failed to fetch endpoints for service %s.%s: %+v", so.spec.Namespace, so.spec.ServiceName, err)
		return
	}

	eh.applyEndpoints(so, eps)
}

func (eh *endpointsHub) applyEndpoints(so *serviceObject, eps *api.Endpoints) {
	var fullEps *pb.ServiceEndpoints
	var observers []*clientObject
	so.WithWLock(func() error {
		so.endpoints = skypbEndpointsToMap(so.spec, eps)
		fullEps = skypbEndpointsToSlice(so.spec, eps)
		observers = so.observers
		return nil
	})

	for _, observer := range observers {
		go func(observer *clientObject, eps *pb.ServiceEndpoints) {
			if len(eps.InstEndpoints) == 0 {
				return
			}

			up := EndpointsUpdate{
				Id:        atomic.AddInt64(&nextUpdateId, 1),
				Endpoints: eps,
			}
			select {
			case <-observer.stopCh:
			case observer.notifyCh <- &up:
			}
		}(observer, fullEps)
	}
}

func (eh *endpointsHub) TrackServiceGraph(req *pb.ResolveRequest, callee *pb.ServiceSpec, callerAddr net.Addr) {
	glog.V(3).Infof("TrackServiceGraph %#v|%#v --> %#v\n", req.CallerServiceId, req.CallerServiceName, callee)

	graphKey := path.Join(prefix.GraphKey, callee.Namespace, callee.ServiceName, "clients", req.CallerServiceName)
	glog.V(5).Infof("etcd set %#v\n", graphKey)

	eh.graphKeysLock.Lock()
	eh.graphKeys[graphKey] = struct{}{}
	eh.graphKeysLock.Unlock()

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	succeeded := false
	for i := 0; i < 3; i++ {
		if _, err := eh.etcdCli.Set(context.Background(), graphKey, timestamp, setGraphOpts); nil != err {
			if i == 2 {
				glog.Warningf("Save service graph key %s in etcd: %#v", graphKey, err)
			}
			continue
		}
		succeeded = true
		break
	}
	if !succeeded {
		glog.V(5).Infof("Failed to save service graph key %s in etcd for 3 times", graphKey)
	}
}

func (eh *endpointsHub) UntrackServiceGraph(req *pb.ResolveRequest, callee *pb.ServiceSpec, callerAddr net.Addr) {
	glog.V(3).Infof("UntrackServiceGraph %#v|%#v --> %#v\n", req.CallerServiceId, req.CallerServiceName, callee)

	graphKey := path.Join(prefix.GraphKey, callee.Namespace, callee.ServiceName, "clients", req.CallerServiceName)

	eh.graphKeysLock.Lock()
	delete(eh.graphKeys, graphKey)
	eh.graphKeysLock.Unlock()
}

func (eh *endpointsHub) startGraphTracking() {
	for range time.Tick(*graphKeyInterval) {
		// Clone the graph keys map.
		keys := make(map[string]struct{})
		eh.graphKeysLock.RLock()
		for k := range eh.graphKeys {
			keys[k] = struct{}{}
		}
		eh.graphKeysLock.RUnlock()

		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		for k := range keys {
			if _, err := eh.etcdCli.Set(context.Background(), k, timestamp, setGraphOpts); nil != err {
				glog.Warningf("Save service graph key %s in etcd err: %#v", k, err)
			}
			// Throttle the traffic to ETCD to 20 keys/sec.
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// Init initializes and returns the endpoint hub.
func Init() EndpointsHub {
	// Start hub only once.
	once.Do(func() {
		hub = &endpointsHub{
			services:      make(serviceMap),
			etcdCli:       CreateEtcdClient(*EtcdEndpoints, true),
			graphKeys:     make(map[string]struct{}),
			graphKeysLock: &sync.RWMutex{},
		}
		prefix.Init(hub.etcdCli)
		if *withinK8s {
			go hub.startK8sWatcher()
		} else {
			go hub.startMainWatcher()
		}
		go hub.startLameDuckWatcher()
		go hub.startGraphTracking()
	})
	return hub
}

// CreateEtcdClient returns a new Etcd client.
func CreateEtcdClient(etcdEndpoints string, required bool) etcd.KeysAPI {
	eps := strings.CsvToSlice(etcdEndpoints)
	if len(eps) == 0 {
		if !required {
			return nil
		}
		glog.Fatalln("Flag --etcd-endpoints is required.")
	}

	glog.Infof("Use etcd endpoints %s.", eps)

	var cli etcd.Client
	for {
		var err error
		if cli, err = etcd.New(etcd.Config{
			Endpoints: eps,
		}); err != nil {
			glog.Errorf("Failed to create etcd client, %v. Will retry after one second.", err)
			time.Sleep(time.Second)
			continue
		}
		if err = cli.Sync(context.Background()); err != nil {
			glog.Errorf("Failed to sync cluster: %v. Will retry after one second.", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	machines := cli.Endpoints()
	if len(machines) == 0 || len(machines[0]) == 0 {
		glog.Fatalln("No etcd machines found")
	}
	glog.Infoln("Found etcd machines:", machines)

	return etcd.NewKeysAPI(cli)
}

func skypbEndpointsToMap(spec *pb.ServiceSpec, eps *api.Endpoints) serviceEndpoints {
	m := make(serviceEndpoints)
	for _, s := range eps.Subsets {
		port := findPort(s.Ports, spec.PortName)
		if port == 0 {
			continue
		}
		for _, addr := range s.Addresses {
			se := ServiceEndpoint{
				IP:   addr.IP,
				Port: port,
			}
			if weight, ok := eps.Labels[calculateWeightKey(addr.IP, port)]; ok {
				if tmpWeight, err := strconv.Atoi(weight); err == nil {
					se.Weight = int32(tmpWeight)
				}
			}
			m[se.toString()] = se
		}
	}
	return m
}

func skypbEndpointsToSlice(spec *pb.ServiceSpec, eps *api.Endpoints) *pb.ServiceEndpoints {
	svcEps := make([]*pb.InstanceEndpoint, 0, len(eps.Subsets))
	for _, s := range eps.Subsets {
		port := findPort(s.Ports, spec.PortName)
		if port == 0 {
			continue
		}
		for _, addr := range s.Addresses {
			ep := pb.InstanceEndpoint{
				Host: addr.IP,
				Port: port,
			}
			if weight, ok := eps.Labels[calculateWeightKey(addr.IP, port)]; ok {
				if tmpWeight, err := strconv.Atoi(weight); err == nil {
					ep.Weight = int32(tmpWeight)
				}
			}
			svcEps = append(svcEps, &ep)
		}
	}

	return &pb.ServiceEndpoints{
		Spec:          spec,
		InstEndpoints: svcEps,
	}
}
