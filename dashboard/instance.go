package dashboard

import (
	"context"
	"fmt"
	"sort"
	"strings"

	etcd "github.com/coreos/etcd/client"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/kataras/iris"

	"github.com/binchencoder/letsgo/service/naming"
	"github.com/binchencoder/skylb-api/lameduck"
	"github.com/binchencoder/skylb/dashboard/db"
	pb "github.com/binchencoder/skylb/dashboard/proto"
	"github.com/binchencoder/ease-gateway/proto/data"
)

var (
	setOpts = etcd.SetOptions{}
	delOpts = etcd.DeleteOptions{}
)

func addInstanceHandler(ctx *iris.Context) {
	req := pb.AddInstanceRequest{}
	resp := pb.AddInstanceResponse{}

	curUser, ok := ctx.Session().Get(sessionUserKey).(*db.User)
	if !ok {
		resp.ErrorMsg = "Not login."
		pbResponse(ctx, &resp)
		return
	}
	if !isAdmin(curUser) {
		resp.ErrorMsg = "No permission."
		pbResponse(ctx, &resp)
		return
	}

	if err := proto.Unmarshal(ctx.PostBody(), &req); err != nil {
		glog.Errorf("Failed to unmarshal AddInstanceRequest, %v", err)
		ctx.Data(iris.StatusBadRequest, nil)
		return
	}

	name, err := naming.ServiceIdToName(data.ServiceId(req.Id))
	if err != nil {
		resp.ErrorMsg = fmt.Sprintf("Service ID %d not found", req.Id)
		pbResponse(ctx, &resp)
		return
	}

	if req.Address == "" || strings.Index(req.Address, ":") == -1 {
		resp.ErrorMsg = fmt.Sprintf("%s is not a valid address", req.Address)
		pbResponse(ctx, &resp)
		return
	}

	parts := strings.Split(req.Address, ":")
	if err := lameduck.SetLameDuckMode(etcdCli, name, lameduck.HostPort(parts[0], parts[1])); err != nil {
		glog.Errorf("Failed to set lameduck for service %s, instance %s, %v", name, req.Address, err)
		pbResponse(ctx, &resp)
		return
	}

	operation := fmt.Sprintf("Add instance %s and put in lameduck mode", req.Address)
	db.CreateLog(curUser.LoginName, req.Id, operation)

	pbResponse(ctx, &resp)
}

func loadInstances(name string, lameducks map[string]string) []*pb.InstanceInfo {
	keyPrefix := "/registry/services/endpoints/default/" + name + "/"
	etcdResp, err := etcdCli.Get(context.Background(), keyPrefix, &etcd.GetOptions{Recursive: true})
	if err != nil {
		glog.Errorf("Failed to load service instances with key prefix %s, %v", keyPrefix, err)
	} else {
		// Note: lameducks will be changed by extractInstances().
		return extractInstances(etcdResp.Node, len(keyPrefix), lameducks)
	}
	return nil
}

func extractInstances(root *etcd.Node, prefixLen int, lameducks map[string]string) []*pb.InstanceInfo {
	instances := []*pb.InstanceInfo{}
	for _, node := range root.Nodes {
		hostAddr := node.Key[prefixLen:]
		hostAddr = strings.Replace(hostAddr, "_", ":", 1)
		isLameduck := false
		if _, ok := lameducks[hostAddr]; ok {
			isLameduck = true
			delete(lameducks, hostAddr)
		}
		instances = append(instances, &pb.InstanceInfo{
			Address:  hostAddr,
			Lameduck: isLameduck,
		})
	}
	for k := range lameducks {
		instances = append(instances, &pb.InstanceInfo{
			Address:  k,
			Lameduck: true,
		})
	}
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].Address < instances[j].Address
	})
	return instances
}
