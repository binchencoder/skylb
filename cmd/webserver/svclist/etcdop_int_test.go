package svclist

import (
	"encoding/json"
	"testing"

	pb "github.com/binchencoder/skylb-api/proto"
	"github.com/binchencoder/skylb/hub"
)

var (
	etcdCli = hub.GetTestEtcdClient()
)

func TestListServices(t *testing.T) {
	list, err := ListServices(etcdCli)
	if nil != err {
		t.Fatal(err)
	}
	//t.Logf("%v", list)

	//b, err := json.Marshal(&list)
	b, err := json.MarshalIndent(&list, "", "\t")
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%s", b)
}

func TestGetDependencies(t *testing.T) {
	buf, err := GetDependencies(etcdCli)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("Dep:\n%s", buf.String())

	buf, err = GetDependencies(etcdCli)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("Dep2:\n%s", buf.String())
}

func TestJson(t *testing.T) {
	ep := hub.ServiceEndpoint{
		IP:   "1.2.3.4",
		Port: 6789,
	}
	b, err := json.MarshalIndent(&ep, "", "\t")
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%s", b)

	ieps := make([]*pb.InstanceEndpoint, 3)
	ieps = append(ieps, &pb.InstanceEndpoint{
		Host: "h1",
		Port: 1234,
	})
	seps := pb.ServiceEndpoints{
		Spec: &pb.ServiceSpec{
			Namespace:   "ns1",
			ServiceName: "sn1",
		},
		InstEndpoints: ieps,
	}
	t.Log(seps)
	b, err = json.MarshalIndent(&seps, "", "\t")
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("%s", b)
}
