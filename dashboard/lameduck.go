package dashboard

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/kataras/iris"

	"github.com/binchencoder/letsgo/service/naming"
	"github.com/binchencoder/skylb-api/lameduck"
	"github.com/binchencoder/skylb/dashboard/db"
	pb "github.com/binchencoder/skylb/dashboard/proto"
	"github.com/binchencoder/ease-gateway/proto/data"
)

func toggleLameduckHandler(ctx *iris.Context) {
	req := pb.ToggleLameduckRequest{}
	resp := pb.ToggleLameduckResponse{}

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
		glog.Errorf("Failed to unmarshal ToggleLameduckRequest, %v", err)
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
	if req.Lameduck {
		if err := lameduck.SetLameDuckMode(etcdCli, name, lameduck.HostPort(parts[0], parts[1])); err != nil {
			glog.Errorf("Failed to set lameduck for service %s, instance %s, %v", name, req.Address, err)
			pbResponse(ctx, &resp)
			return
		}
	} else {
		if err := lameduck.UnsetLameDuckMode(etcdCli, name, lameduck.HostPort(parts[0], parts[1])); err != nil {
			glog.Errorf("Failed to unset lameduck for service %s, instance %s, %v", name, req.Address, err)
			pbResponse(ctx, &resp)
			return
		}
	}

	var operation string
	if req.Lameduck {
		operation = fmt.Sprintf("Put %s in lameduck mode", req.Address)
	} else {
		operation = fmt.Sprintf("Take %s out of lameduck mode", req.Address)
	}
	db.CreateLog(curUser.LoginName, req.Id, operation)

	pbResponse(ctx, &resp)
}
