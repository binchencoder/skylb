package dashboard

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/iris"
	"golang.org/x/net/context"

	"github.com/binchencoder/skylb/dashboard/db"
	pb "github.com/binchencoder/skylb/dashboard/proto"
)

func getLogsHandler(ctx *iris.Context) {
	req := pb.GetLogsRequest{}
	resp := pb.GetLogsResponse{}

	if err := proto.Unmarshal(ctx.PostBody(), &req); err != nil {
		glog.Errorf("Failed to unmarshal GetLogsRequest, %v", err)
		ctx.Data(iris.StatusBadRequest, nil)
		return
	}

	_, ok := ctx.Session().Get(sessionUserKey).(*db.User)
	if !ok {
		resp.ErrorMsg = "Not login."
		pbResponse(ctx, &resp)
		return
	}

	logs, err := db.GetLogs(context.Background(), req.Operator, req.ServiceId)
	if err != nil {
		resp.ErrorMsg = err.Error()
		pbResponse(ctx, &resp)
		return
	}

	resp.Operator = req.Operator
	resp.ServiceId = req.ServiceId
	for _, l := range logs {
		resp.Logs = append(resp.Logs, l.ToLogInfo())
	}
	pbResponse(ctx, &resp)
}
