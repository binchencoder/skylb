package hub

import (
	"net"
	"testing"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	pb "binchencoder.com/skylb-api/proto"
)

func TestTrackServiceGraph(t *testing.T) {
	etcdCli := GetTestEtcdClient()
	eh := endpointsHub{
		etcdCli: etcdCli,
	}
	req := &pb.ResolveRequest{
		CallerServiceId:   1,
		CallerServiceName: "caller1",
	}
	callee := &pb.ServiceSpec{
		Namespace:   "default",
		ServiceName: "called1",
		PortName:    "grpc",
	}
	addr, _ := net.ResolveIPAddr("ip4", "1.2.3.4")

	// Normal case
	eh.TrackServiceGraph(req, callee, addr)

	getOpts := &etcd.GetOptions{}
	if _, err := etcdCli.Get(context.Background(), "/skylb/graph/default/called1/clients/caller1", getOpts); nil != err {
		t.Errorf("TrackServiceGraph failed:%v", err)
	}

	// Now simulate existing wrong key ".../clients"
	req.CallerServiceName = "caller2"
	callee.ServiceName = "called2"
	delOpts := &etcd.DeleteOptions{
		Recursive: true,
	}
	if _, err := etcdCli.Delete(context.Background(), "/skylb/graph/default/called2/clients", delOpts); nil != err {
		t.Errorf("Init test data err:%v", err)
	}
	if _, err := etcdCli.Set(context.Background(), "/skylb/graph/default/called2/clients", "stub", setGraphOpts); nil != err {
		t.Errorf("Simulate test data err:%v", err)
	}
	eh.TrackServiceGraph(req, callee, addr)
	if _, err := etcdCli.Get(context.Background(), "/skylb/graph/default/called2/clients/caller2", getOpts); nil != err {
		t.Errorf("TrackServiceGraph failed:%v", err)
	}
	// addr will be saved since second call, not quite good, but acceptable.
	eh.TrackServiceGraph(req, callee, addr)
	if _, err := etcdCli.Get(context.Background(), "/skylb/graph/default/called2/clients/caller2", getOpts); nil != err {
		t.Errorf("TrackServiceGraph failed:%v", err)
	}
}

func TestFetchEndpoints(t *testing.T) {
	etcdCli := GetTestEtcdClient()
	eh := endpointsHub{
		etcdCli: etcdCli,
	}
	eps, err := eh.fetchEndpoints("default", "shared-test-server-service")
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%v", eps)
}
