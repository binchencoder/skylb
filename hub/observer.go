package hub

import (
	"flag"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	prom "github.com/prometheus/client_golang/prometheus"
	api "k8s.io/api/core/v1"

	pb "github.com/binchencoder/skylb-api/proto"
)

const (
	ChanCapMultiplication = 10
	labelHttpStatus       = "status_code"
)

var (
	autoRectifyInterval = flag.Duration("auto-rectify-interval", 60*time.Second, "The interval of auto rectification.")

	addObserverGauge = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "add_observer_gauge",
			Help:      "SkyLB add observer gauge.",
		},
		[]string{"service"},
	)
	removeObserverGauge = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "infra",
			Subsystem: "skylb",
			Name:      "remove_observer_gauge",
			Help:      "SkyLB remove observer gauge.",
		},
		[]string{"service"},
	)
)

func init() {
	prom.MustRegister(addObserverGauge)
	prom.MustRegister(removeObserverGauge)
}

// AddObserver adds an observer of the given service specs for the given
// clientAddr. When service endpoints changed, it notifies the observer
// through the returned channel.
func (eh *endpointsHub) AddObserver(specs []*pb.ServiceSpec, clientAddr string, resolveFull bool) (<-chan *EndpointsUpdate, error) {
	notifyCh := make(chan *EndpointsUpdate, ChanCapMultiplication*len(specs))

	for _, spec := range specs {
		glog.V(2).Infof("Resolve service %s.%s on port name %q from client %s", spec.Namespace, spec.ServiceName, spec.PortName, clientAddr)
		label := fmt.Sprintf("%s.%s", spec.Namespace, spec.ServiceName)
		addObserverGauge.WithLabelValues(label).Inc()

		co := &clientObject{
			spec:        spec,
			clientAddr:  clientAddr,
			notifyCh:    notifyCh,
			stopCh:      make(chan struct{}),
			resolveFull: resolveFull,
		}

		var eps *api.Endpoints
		var err error
		if *withinK8s {
			eps, err = eh.fetchK8sEndpoints(spec.Namespace, spec.ServiceName)
		} else {
			eps, err = eh.fetchEndpoints(spec.Namespace, spec.ServiceName)
		}
		if err != nil {
			return nil, err
		}

		key := eh.calculateKey(spec.Namespace, spec.ServiceName)
		var so *serviceObject
		err = eh.WithWLock(func() error {
			glog.V(3).Infof("Received initial endpoints for client %s: %+v.", clientAddr, eps)

			epsMap := skypbEndpointsToMap(spec, eps)
			var ok bool
			if so, ok = eh.services[key]; !ok {
				so = &serviceObject{
					spec:      spec,
					endpoints: epsMap,
				}
				eh.services[key] = so

				if !*withinK8s {
					// Periodically update the endpoints so that client gets a
					// chance to rectify its endpoint list.
					go func() {
						ticker := time.NewTicker(*autoRectifyInterval)
						for range ticker.C {
							glog.V(3).Infof("Automatic endpoints rectification for %s.", key)
							eh.updateEndpoints(key)
						}
					}()
				}
			}

			up := EndpointsUpdate{
				Id:        atomic.AddInt64(&nextUpdateId, 1),
				Endpoints: diffEndpoints(spec, nil, epsMap),
			}
			notifyCh <- &up
			return nil
		})
		if err != nil {
			close(notifyCh)
			return nil, err
		}

		so.WithWLock(func() error {
			so.observers = append(so.observers, co)
			return nil
		})
	}

	return notifyCh, nil
}

// RemoveObserver removes the observer for the given service specs for the
// given clientAddr.
func (eh *endpointsHub) RemoveObserver(specs []*pb.ServiceSpec, clientAddr string) {
	for _, spec := range specs {
		label := fmt.Sprintf("%s.%s", spec.Namespace, spec.ServiceName)
		removeObserverGauge.WithLabelValues(label).Inc()

		key := eh.calculateKey(spec.Namespace, spec.ServiceName)

		var so *serviceObject
		var ok bool
		eh.WithRLock(func() error {
			so, ok = eh.services[key]
			return nil
		})
		if !ok {
			return
		}

		so.WithWLock(func() error {
			so.observers = removeObserverFromSlice(so.observers, spec, clientAddr)
			return nil
		})
	}
}

func removeObserverFromSlice(obs []*clientObject, spec *pb.ServiceSpec, clientAddr string) []*clientObject {
	remaining := make([]*clientObject, 0, len(obs))
	for _, ob := range obs {
		if ob.clientAddr == clientAddr && ob.spec.String() == spec.String() {
			// Stop all goroutines for this observer.
			close(ob.stopCh)
		} else {
			remaining = append(remaining, ob)
		}
	}
	return remaining
}
