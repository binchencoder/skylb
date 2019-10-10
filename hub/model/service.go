package model

import (
	"flag"
	"sync"
	"time"

	"github.com/golang/glog"

	pb "binchencoder.com/skylb-api/proto"
	"binchencoder.com/skylb/hub/util"
)

var (
	autoRectifyInterval = flag.Duration("auto-rectify-interval", 60*time.Second, "The interval of auto rectification.")
)

// ServiceObject defines the interface to manages service endpoints
// and client observers for one service.
type ServiceObject interface {
	// Spec returns the spec of the service.
	Spec() *pb.ServiceSpec

	// AddObserver adds a service client observer to the service object.
	AddObserver(co ClientObserver)

	// RemoveObservers removes all client observers with the specified clientAddr
	// from the service object.
	RemoveObservers(clientAddr string)

	// SetEndpoints sets the service enpoints.
	SetEndpoints(endpoints *pb.ServiceEndpoints)
}

// serviceObject implements interface ServiceObject.
type serviceObject struct {
	lock sync.Mutex

	stopCh chan struct{} // Channel to stop the ticker.

	spec      *pb.ServiceSpec
	key       string
	endpoints *pb.ServiceEndpoints
	observers []ClientObserver

	listener func(key string) // The auto endpoints rectification listener.
}

func (so *serviceObject) Spec() *pb.ServiceSpec {
	return so.spec
}

func (so *serviceObject) AddObserver(co ClientObserver) {
	so.lock.Lock()
	defer so.lock.Unlock()

	glog.V(4).Infof("Adding observer %s for service %s.%s.\n", co.ClientAddr(), so.spec.GetNamespace(), so.spec.GetServiceName())

	if len(so.observers) == 0 {
		// Start the auto rectify goroutine.
		so.stopCh = make(chan struct{})
		go so.startAutoRectify()
	}

	so.observers = append(so.observers, co)

	if so.endpoints != nil && len(so.endpoints.InstEndpoints) > 0 {
		// Notify client observer so that it will not block
		// the current goroutine.
		go co.Notify(so.endpoints)
	}
}

func (so *serviceObject) RemoveObservers(clientAddr string) {
	so.lock.Lock()
	defer so.lock.Unlock()

	if len(so.observers) == 0 {
		return
	}

	remaining := make([]ClientObserver, 0, len(so.observers))
	for _, observer := range so.observers {
		if observer.ClientAddr() == clientAddr && observer.Spec().String() == so.spec.String() {
			observer.Close()
		} else {
			remaining = append(remaining, observer)
		}
	}
	so.observers = remaining

	if len(so.observers) == 0 && so.stopCh != nil {
		// Stop the auto rectify goroutine.
		close(so.stopCh)
		so.stopCh = nil
	}
}

func (so *serviceObject) SetEndpoints(endpoints *pb.ServiceEndpoints) {
	so.lock.Lock()
	defer so.lock.Unlock()

	glog.V(4).Infof("Setting endpoints %s for service %s.%s.\n", endpoints, so.spec.GetNamespace(), so.spec.GetServiceName())

	so.endpoints = endpoints
	for _, co := range so.observers {
		go co.Notify(endpoints)
	}
}

// startAutoRectify periodically updates the endpoints so that client gets a
// chance to rectify its endpoint list.
func (so *serviceObject) startAutoRectify() {
	ticker := time.NewTicker(*autoRectifyInterval)
	for {
		select {
		case <-so.stopCh:
			// Exit to avoid goroutine leak.
			ticker.Stop()
			return
		case <-ticker.C:
			glog.V(3).Infof("Automatic endpoints rectification for %s.", so.key)
			if so.listener != nil {
				so.listener(so.key)
			}
		}
	}
}

// NewServiceObject creates and returns a new ServiceObject.
func NewServiceObject(spec *pb.ServiceSpec, listener func(key string)) ServiceObject {
	so := serviceObject{
		spec:      spec,
		key:       util.CalculateKey(spec.Namespace, spec.ServiceName),
		endpoints: &pb.ServiceEndpoints{},
		observers: []ClientObserver{},
		listener:  listener,
	}

	return &so
}
