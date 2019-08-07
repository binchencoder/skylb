package svclist

import (
	"bytes"
	"encoding/json"
	"path"
	"strconv"
	"strings"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	api "k8s.io/api/core/v1"

	"github.com/binchencoder/skylb-api/prefix"
	pb "github.com/binchencoder/skylb-api/proto"
	"github.com/binchencoder/skylb/hub"
)

var (
	getOpts = &etcd.GetOptions{
		Recursive: true,
		Sort:      true,
	}
)

func ListServices(etcdCli etcd.KeysAPI) ([]*pb.ServiceEndpoints, error) {
	sepsList := make([]*pb.ServiceEndpoints, 0, 200)
	resp, err := etcdCli.Get(context.Background(), prefix.EndpointsKey, getOpts)
	if nil != err {
		return nil, err
	}
	logLevel := glog.Level(5)
	for _, ns := range resp.Node.Nodes {
		glog.V(logLevel).Infof("namespace %#v", ns.Key)
		for _, svc := range ns.Nodes {
			glog.V(logLevel).Infof("> svc %#v", svc.Key)
			seps := &pb.ServiceEndpoints{
				Spec: &pb.ServiceSpec{
					Namespace:   path.Base(ns.Key),
					ServiceName: path.Base(svc.Key),
					//PortName: // TODO(fuyc): fill it if need to.
				},
			}
			ieps := make([]*pb.InstanceEndpoint, 0, 30)
			for _, si := range svc.Nodes {
				glog.V(logLevel).Infof(">> inst %#v", si.Key)
				hostPort := path.Base(si.Key)
				ss := strings.Split(hostPort, "_")
				if len(ss) < 2 {
					glog.Warningf("Invalid hostPort: %s", hostPort)
					continue
				}
				port, err := strconv.Atoi(ss[1])
				if err != nil {
					glog.Warningf("Invalid port: %s %v", ss[1], err)
					continue
				}
				iep := &pb.InstanceEndpoint{
					Host: ss[0],
					Port: int32(port),
				}
				ieps = append(ieps, iep)
				if seps.Spec.PortName == "" {
					eps := &api.Endpoints{}
					err = json.Unmarshal([]byte(si.Value), eps)
					if err != nil {
						glog.V(logLevel).Infof("unmarshaling value: %#v but got error", err, ns.Value)
						continue
					}
					if len(eps.Subsets) > 0 && len(eps.Subsets[0].Ports) > 0 {
						seps.Spec.PortName = eps.Subsets[0].Ports[0].Name
					}
				}
			}
			seps.InstEndpoints = ieps
			sepsList = append(sepsList, seps)
		}
	}
	return sepsList, nil
}

type Clients struct {
	Seps []*pb.ServiceEndpoints
}

func GetDependencies(etcdCli etcd.KeysAPI) (bytes.Buffer, error) {
	var buf bytes.Buffer
	if err := hub.BuildDependencies(etcdCli); nil != err {
		return buf, err
	}
	// Traverse roots.
	rs := hub.FindRoots()
	glog.V(3).Infof("Roots %v", rs)
	for _, r := range rs {
		traverseCalls(&buf, hub.SimpleCallingMap[r][0].Caller, 0)
	}
	return buf, nil
}

func traverseCalls(buf *bytes.Buffer, caller string, level int) {
	for i := 0; i < level; i++ {
		buf.WriteString("&nbsp;&nbsp;&nbsp;&nbsp;")
	}
	if level > 0 {
		buf.WriteString("|_")
	}
	buf.WriteString(caller)
	buf.WriteString("<br/>\n")

	// Efficient lookup.
	cps := hub.SimpleCallingMap[caller]
	for _, cp := range cps {
		if cp.Callee == caller {
			glog.Warningf("Loop detected: %s", caller)
			continue
		}
		traverseCalls(buf, cp.Callee, level+1)
	}
}
