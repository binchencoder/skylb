package model

import (
	"flag"
	"sync"
	"time"

	pb "github.com/binchencoder/skylb-api/proto"
)

var (
	observerNotifyTimeout = flag.Duration("observer-notify-timeout", 5*time.Second, "The timeout to send endpoint updates notification to client.")
)

// ClientObserver defines the interface of a client observer for a gRPC service.
type ClientObserver interface {
	// ClientAddr returns the client address of a client observer.
	ClientAddr() string

	// Spec returns the spec of the service.
	Spec() *pb.ServiceSpec

	// Notify notifies the client observer of the given service endpoints.
	Notify(eps *pb.ServiceEndpoints)

	// Close closes the client observer.
	Close()
}

// clientObserver implements interface ClientObserver.
type clientObserver struct {
	lock sync.Mutex

	spec       *pb.ServiceSpec
	clientAddr string
	closed     bool
	notifyCh   chan<- *pb.ServiceEndpoints
}

func (co *clientObserver) ClientAddr() string {
	return co.clientAddr
}

func (co *clientObserver) Spec() *pb.ServiceSpec {
	return co.spec
}

func (co *clientObserver) Notify(eps *pb.ServiceEndpoints) {
	co.lock.Lock()
	defer co.lock.Unlock()

	if co.closed {
		return
	}

	timer := time.NewTimer(*observerNotifyTimeout)
	select {
	case <-timer.C:
		return
	case co.notifyCh <- eps:
		if !timer.Stop() {
			<-timer.C
		}
	}
}

func (co *clientObserver) Close() {
	co.lock.Lock()
	defer co.lock.Unlock()

	co.closed = true
}

// NewClientObserver creates and returns a new ClientObserver.
func NewClientObserver(spec *pb.ServiceSpec, clientAddr string, notifyCh chan<- *pb.ServiceEndpoints) ClientObserver {
	return &clientObserver{
		spec:       spec,
		clientAddr: clientAddr,
		notifyCh:   notifyCh,
	}
}
