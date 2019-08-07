package rpc

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"errors"

	pb "github.com/binchencoder/skylb-api/proto"
	"github.com/binchencoder/skylb/hub"
	data "github.com/binchencoder/ease-gateway/proto/data"
)

// EndpointsHubMock mocks interface EndpointsHubMock.
type EndpointsHubMock struct {
	mock.Mock
}

func (ephm *EndpointsHubMock) AddObserver(specs []*pb.ServiceSpec, clientAddr string, resolveFull bool) (<-chan *hub.EndpointsUpdate, error) {
	args := ephm.Called(specs, clientAddr, resolveFull)
	if res, ok := args.Get(0).(chan *hub.EndpointsUpdate); ok {
		return res, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ephm *EndpointsHubMock) RemoveObserver(specs []*pb.ServiceSpec, clientAddr string) {
	ephm.Called(specs, clientAddr)
}

func (ephm *EndpointsHubMock) InsertEndpoint(spec *pb.ServiceSpec, host string, port, weight int32) error {
	args := ephm.Called(spec, host, port, weight)
	return args.Error(0)
}

func (ephm *EndpointsHubMock) UpsertEndpoint(spec *pb.ServiceSpec, host string, port, weight int32) error {
	args := ephm.Called(spec, host, port, weight)
	return args.Error(0)
}

func (ephm *EndpointsHubMock) TrackServiceGraph(req *pb.ResolveRequest, callee *pb.ServiceSpec, callerAddr net.Addr) {
	ephm.Called(req, callee, callerAddr)
}

func (ephm *EndpointsHubMock) UntrackServiceGraph(req *pb.ResolveRequest, callee *pb.ServiceSpec, callerAddr net.Addr) {
	ephm.Called(req, callee, callerAddr)
}

// ResolveServer mocks interface pb.Skylb_ResolveServer.
type ResolveServer struct {
	mock.Mock
}

func (rs *ResolveServer) Context() context.Context {
	args := rs.Called()
	if res, ok := args.Get(0).(context.Context); ok {
		return res
	}
	return nil
}

func (rs *ResolveServer) SendMsg(m interface{}) error {
	args := rs.Called(m)
	return args.Error(0)
}

func (rs *ResolveServer) RecvMsg(m interface{}) error {
	args := rs.Called(m)
	return args.Error(0)
}

func (rs *ResolveServer) SetHeader(md metadata.MD) error {
	args := rs.Called(md)
	return args.Error(0)
}

func (rs *ResolveServer) SetTrailer(md metadata.MD) {
	rs.Called(md)
}

func (rs *ResolveServer) SendHeader(md metadata.MD) error {
	args := rs.Called(md)
	return args.Error(0)
}

func (rs *ResolveServer) Send(resp *pb.ResolveResponse) error {
	args := rs.Called(resp)
	return args.Error(0)
}

// ReportLoadServer mocks interface pb.Skylb_ReportLoadServer.
type ReportLoadServer struct {
	mock.Mock
}

func (rls *ReportLoadServer) Context() context.Context {
	args := rls.Called()
	if res, ok := args.Get(0).(context.Context); ok {
		return res
	}
	return nil
}

func (rls *ReportLoadServer) SendMsg(m interface{}) error {
	args := rls.Called(m)
	return args.Error(0)
}

func (rls *ReportLoadServer) RecvMsg(m interface{}) error {
	args := rls.Called(m)
	return args.Error(0)
}
func (rls *ReportLoadServer) SetHeader(md metadata.MD) error {
	args := rls.Called(md)
	return args.Error(0)
}

func (rls *ReportLoadServer) SetTrailer(md metadata.MD) {
	rls.Called(md)
}

func (rls *ReportLoadServer) SendHeader(md metadata.MD) error {
	args := rls.Called(md)
	return args.Error(0)
}

func (rls *ReportLoadServer) SendAndClose(resp *pb.ReportLoadResponse) error {
	args := rls.Called(resp)
	return args.Error(0)
}

func (rls *ReportLoadServer) Recv() (*pb.ReportLoadRequest, error) {
	args := rls.Called()
	if res, ok := args.Get(0).(*pb.ReportLoadRequest); ok {
		return res, nil
	}
	return nil, args.Error(1)
}

func TestOpToString(t *testing.T) {
	for op, expected := range map[pb.Operation]string{
		pb.Operation_Add:    "ADD",
		pb.Operation_Delete: "DELETE",
		pb.Operation(10000): "",
	} {
		str := opToString(op)
		if str != expected {
			t.Errorf("expect %s but got %s", expected, str)
		}
	}
}
func TestResolve(t *testing.T) {
	spec := pb.ServiceSpec{
		Namespace:   "default",
		PortName:    "grpc",
		ServiceName: "test-service",
	}
	ep1 := pb.InstanceEndpoint{
		Op:   pb.Operation_Add,
		Host: "172.0.0.101",
		Port: 8080,
	}
	ep2 := pb.InstanceEndpoint{
		Op:   pb.Operation_Add,
		Host: "172.0.0.101",
		Port: 8080,
	}

	addr, _ := net.ResolveIPAddr("ip", "192.168.0.101")
	ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: addr})

	resp := pb.ResolveResponse{
		SvcEndpoints: &pb.ServiceEndpoints{
			Spec:          &spec,
			InstEndpoints: []*pb.InstanceEndpoint{&ep1, &ep2},
		},
	}
	stream := new(ResolveServer)
	stream.On("Context").Return(ctx)
	stream.On("Send", &resp).Return(nil)

	eh := new(EndpointsHubMock)
	s := &skylbServer{
		epsHub: eh,
	}

	req := pb.ResolveRequest{
		CallerServiceId:      data.ServiceId_SHARED_TEST_CLIENT_SERVICE,
		CallerServiceName:    data.ServiceId_SHARED_TEST_CLIENT_SERVICE.String(),
		ResolveFullEndpoints: true,
		Services:             []*pb.ServiceSpec{&spec},
	}

	ch := make(chan *hub.EndpointsUpdate, 1)
	ch <- &hub.EndpointsUpdate{
		Id: 100,
		Endpoints: &pb.ServiceEndpoints{
			Spec:          &spec,
			InstEndpoints: []*pb.InstanceEndpoint{&ep1, &ep2},
		},
	}
	close(ch)
	eh.On("TrackServiceGraph", &req, &spec, addr)
	eh.On("UntrackServiceGraph", &req, &spec, addr)
	eh.On("AddObserver", []*pb.ServiceSpec{&spec}, "192.168.0.101", true).Return(ch, nil)
	eh.On("RemoveObserver", []*pb.ServiceSpec{&spec}, "192.168.0.101")

	err := s.Resolve(&req, stream)
	if err != nil {
		t.Errorf("expect no error but got %v", err)
	}
}

func TestResolve_timeout(t *testing.T) {
	spec := pb.ServiceSpec{
		Namespace:   "default",
		PortName:    "grpc",
		ServiceName: "test-service",
	}
	ep1 := pb.InstanceEndpoint{
		Op:   pb.Operation_Add,
		Host: "172.0.0.101",
		Port: 8080,
	}
	ep2 := pb.InstanceEndpoint{
		Op:   pb.Operation_Add,
		Host: "172.0.0.101",
		Port: 8080,
	}

	addr, _ := net.ResolveIPAddr("ip", "192.168.0.101")
	ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: addr})

	resp := pb.ResolveResponse{
		SvcEndpoints: &pb.ServiceEndpoints{
			Spec:          &spec,
			InstEndpoints: []*pb.InstanceEndpoint{&ep1, &ep2},
		},
	}
	stream := new(ResolveServer)
	stream.On("Context").Return(ctx)
	// Force sending to timeout.
	stream.On("Send", &resp).WaitUntil(time.After(5 * time.Second)).Return(nil)

	eh := new(EndpointsHubMock)
	s := &skylbServer{
		epsHub: eh,
	}

	req := pb.ResolveRequest{
		CallerServiceId:      data.ServiceId_SHARED_TEST_CLIENT_SERVICE,
		CallerServiceName:    data.ServiceId_SHARED_TEST_CLIENT_SERVICE.String(),
		ResolveFullEndpoints: true,
		Services:             []*pb.ServiceSpec{&spec},
	}

	ch := make(chan *hub.EndpointsUpdate, 10)
	ch <- &hub.EndpointsUpdate{
		Id: 100,
		Endpoints: &pb.ServiceEndpoints{
			Spec:          &spec,
			InstEndpoints: []*pb.InstanceEndpoint{&ep1, &ep2},
		},
	}
	eh.On("TrackServiceGraph", &req, &spec, addr)
	eh.On("UntrackServiceGraph", &req, &spec, addr)
	eh.On("AddObserver", []*pb.ServiceSpec{&spec}, "192.168.0.101", true).Return(ch, nil)
	eh.On("RemoveObserver", []*pb.ServiceSpec{&spec}, "192.168.0.101")

	// Set a short timeout.
	*flagNotifyTimeout = time.Second

	// Note that when timeout, the stream will be discarded, so although ch was
	// not closed here, the test can still terminate.

	err := s.Resolve(&req, stream)
	if err == nil {
		t.Errorf("expect non-nil error")
	}
}

func TestResolve_errToSend(t *testing.T) {
	spec := pb.ServiceSpec{
		Namespace:   "default",
		PortName:    "grpc",
		ServiceName: "test-service",
	}
	ep1 := pb.InstanceEndpoint{
		Op:   pb.Operation_Add,
		Host: "172.0.0.101",
		Port: 8080,
	}
	ep2 := pb.InstanceEndpoint{
		Op:   pb.Operation_Add,
		Host: "172.0.0.101",
		Port: 8080,
	}
	sendErr := errors.New("failed to send")

	addr, _ := net.ResolveIPAddr("ip", "192.168.0.101")
	ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: addr})

	resp := pb.ResolveResponse{
		SvcEndpoints: &pb.ServiceEndpoints{
			Spec:          &spec,
			InstEndpoints: []*pb.InstanceEndpoint{&ep1, &ep2},
		},
	}
	stream := new(ResolveServer)
	stream.On("Context").Return(ctx)
	stream.On("Send", &resp).Return(sendErr)

	eh := new(EndpointsHubMock)
	s := &skylbServer{
		epsHub: eh,
	}

	req := pb.ResolveRequest{
		CallerServiceId:      data.ServiceId_SHARED_TEST_CLIENT_SERVICE,
		CallerServiceName:    data.ServiceId_SHARED_TEST_CLIENT_SERVICE.String(),
		ResolveFullEndpoints: true,
		Services:             []*pb.ServiceSpec{&spec},
	}

	ch := make(chan *hub.EndpointsUpdate, 10)
	ch <- &hub.EndpointsUpdate{
		Id: 100,
		Endpoints: &pb.ServiceEndpoints{
			Spec:          &spec,
			InstEndpoints: []*pb.InstanceEndpoint{&ep1, &ep2},
		},
	}
	eh.On("TrackServiceGraph", &req, &spec, addr)
	eh.On("UntrackServiceGraph", &req, &spec, addr)
	eh.On("AddObserver", []*pb.ServiceSpec{&spec}, "192.168.0.101", true).Return(ch, nil)
	eh.On("RemoveObserver", []*pb.ServiceSpec{&spec}, "192.168.0.101")

	// Note that when error to send, the stream should be discarded,
	// so although ch was not closed here, the test can still terminate.

	if err := s.Resolve(&req, stream); err == nil {
		t.Errorf("expect non-nil error")
	}
}

func TestReportLoad(t *testing.T) {
	spec := pb.ServiceSpec{
		Namespace:   "default",
		PortName:    "grpc",
		ServiceName: "test-service",
	}

	addr, _ := net.ResolveIPAddr("ip", "192.168.0.101")
	ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: addr})

	req := pb.ReportLoadRequest{
		Spec: &spec,
		Port: 8000,
	}

	quitErr := errors.New("quit")

	stream := ReportLoadServer{}
	stream.On("Context").Return(ctx)
	stream.On("Recv").Once().Return(&req, nil)
	stream.On("Recv").Once().Return(nil, quitErr)

	eh := new(EndpointsHubMock)
	eh.On("UpsertEndpoint", req.Spec, "192.168.0.101", req.Port, int32(0)).Return(nil)
	eh.On("InsertEndpoint", req.Spec, "192.168.0.101", req.Port, int32(0)).Return(nil)

	s := &skylbServer{
		epsHub: eh,
	}

	if err := s.ReportLoad(&stream); err == nil {
		t.Errorf("expect non-nil error")
	}
}
