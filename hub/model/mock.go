package model

import (
	"github.com/stretchr/testify/mock"

	pb "github.com/binchencoder/skylb-api/proto"
)

// ClientObserverMock mocks ClientObserver.
type ClientObserverMock struct {
	mock.Mock
}

func (com *ClientObserverMock) Spec() *pb.ServiceSpec {
	args := com.Called()
	if res, ok := args.Get(0).(*pb.ServiceSpec); ok {
		return res
	}
	return nil
}

func (com *ClientObserverMock) ClientAddr() string {
	args := com.Called()
	return args.String(0)
}

func (com *ClientObserverMock) Notify(eps *pb.ServiceEndpoints) {
	com.Called(eps)
}

func (com *ClientObserverMock) Close() {
	com.Called()
}

// ServiceObjectMock mocks ServiceObject.
type ServiceObjectMock struct {
	mock.Mock
}

func (so *ServiceObjectMock) Spec() *pb.ServiceSpec {
	args := so.Called()
	if res, ok := args.Get(0).(*pb.ServiceSpec); ok {
		return res
	}
	return nil
}

func (so *ServiceObjectMock) AddObserver(co ClientObserver) {
	so.Called(co)
}

func (so *ServiceObjectMock) RemoveObservers(clientAddr string) {
	so.Called(clientAddr)
}

func (so *ServiceObjectMock) SetEndpoints(endpoints *pb.ServiceEndpoints) {
	so.Called(endpoints)
}
