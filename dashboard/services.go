package dashboard

import (
	"context"
	"sort"
	"strings"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/iris"

	"binchencoder.com/skylb-api/lameduck"
	pb "binchencoder.com/skylb/dashboard/proto"
	"binchencoder.com/skylb/dashboard/util"
	vex "binchencoder.com/gateway-proto/data"
)

func getAllServicesHandler(ctx *iris.Context) {
	req := pb.GetAllServicesRequest{}
	resp := pb.GetAllServicesResponse{}

	if err := proto.Unmarshal(ctx.PostBody(), &req); err != nil {
		glog.Errorf("Failed to unmarshal GetUsersRequest, %v", err)
		ctx.Data(iris.StatusBadRequest, nil)
		return
	}

	for k, v := range vex.ServiceId_name {
		if v == "_" {
			continue
		}

		si := pb.ServiceInfo{
			Id:   k,
			Name: strings.ToLower(strings.Replace(v, "_", "-", -1)),
		}
		resp.Services = append(resp.Services, &si)
	}

	sort.Slice(resp.Services, func(i, j int) bool {
		return resp.Services[i].Name < resp.Services[j].Name
	})

	pbResponse(ctx, &resp)
}

func getServiceByIdHandler(ctx *iris.Context) {
	req := pb.GetServiceByIdRequest{}
	resp := pb.GetServiceByIdResponse{}

	if err := proto.Unmarshal(ctx.PostBody(), &req); err != nil {
		glog.Errorf("Failed to unmarshal GetUsersRequest, %v", err)
		ctx.Data(iris.StatusBadRequest, nil)
		return
	}

	resp.Service = &pb.ServiceInfo{
		Id:   req.Id,
		Name: "Unknown",
	}
	if name, ok := util.ServiceIDsToNames[req.Id]; ok {
		resp.Service.Name = name
	}

	ctxt := context.Background()

	if resp.Service.Name == "" {
		glog.Errorf("Unknown name for service ID %d", req.Id)
		pbResponse(ctx, &resp)
		return
	}

	resp.Service.Instances = loadInstances(resp.Service.Name, lameduck.LoadLameducks(etcdCli, resp.Service.Name))

	keyPrefix := "/skylb/graph/default/"
	etcdResp, err := etcdCli.Get(ctxt, keyPrefix, &etcd.GetOptions{Recursive: true})
	if err != nil {
		glog.Errorf("Failed to load service graph with key prefix %s, %v", keyPrefix, err)
	} else {
		resp.Service.Incomings = extractIncomings(etcdResp.Node, util.ServiceNamesToIds, resp.Service.Name)
		resp.Service.Outgoings = extractOutgoings(etcdResp.Node, util.ServiceNamesToIds, resp.Service.Name)
	}

	pbResponse(ctx, &resp)
}

func extractIncomings(root *etcd.Node, names2ids map[string]int32, name string) []*pb.ServiceInfo {
	incomings := make([]*pb.ServiceInfo, 0, 10)
	rootkey := "/skylb/graph/default/" + name
	cliKey := "/skylb/graph/default/" + name + "/clients/"
	for _, node := range root.Nodes {
		if node.Key != rootkey {
			continue
		}
		var cliNode *etcd.Node
		for _, n := range node.Nodes {
			if strings.HasSuffix(n.Key, "/clients") {
				cliNode = n
				break
			}
		}
		if cliNode == nil {
			continue
		}
		for _, n := range cliNode.Nodes {
			cliName := n.Key[len(cliKey):]
			if cliName == "" || cliName == "addr" {
				continue
			}
			if cliName == name {
				// Ignore when service call itself.
				continue
			}
			if _, ok := names2ids[cliName]; !ok {
				glog.Errorf("Found unknown service name in dependency graph: %s.", cliName)
			} else {
				incomings = append(incomings, &pb.ServiceInfo{
					Id:   names2ids[cliName],
					Name: cliName,
				})
			}
		}
	}
	return incomings
}

func extractOutgoings(root *etcd.Node, names2ids map[string]int32, name string) []*pb.ServiceInfo {
	outgoings := make([]*pb.ServiceInfo, 0, 10)
	rootkey := "/skylb/graph/default/"
	for _, node := range root.Nodes {
		svcName := node.Key[len(rootkey):]
		if svcName == name {
			continue
		}
		var cliNode *etcd.Node
		for _, n := range node.Nodes {
			if strings.HasSuffix(n.Key, "/clients") {
				cliNode = n
				break
			}
		}
		if cliNode == nil {
			continue
		}
		cliKey := rootkey + svcName + "/clients/"
		for _, n := range cliNode.Nodes {
			cliName := n.Key[len(cliKey):]
			if cliName == "" || cliName == "addr" {
				continue
			}
			if cliName != name {
				continue
			}
			if _, ok := names2ids[svcName]; !ok {
				glog.Errorf("Found unknown service name in dependency graph: %s.", svcName)
			} else {
				outgoings = append(outgoings, &pb.ServiceInfo{
					Id:   names2ids[svcName],
					Name: svcName,
				})
			}
		}
	}
	return outgoings
}
