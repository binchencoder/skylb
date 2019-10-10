package dashboard

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/iris"
	"golang.org/x/net/context"

	"binchencoder.com/skylb/dashboard/db"
	pb "binchencoder.com/skylb/dashboard/proto"
)

func getUsersHandler(ctx *iris.Context) {
	req := pb.GetUsersRequest{}
	resp := pb.GetUsersResponse{}

	if err := proto.Unmarshal(ctx.PostBody(), &req); err != nil {
		glog.Errorf("Failed to unmarshal GetUsersRequest, %v", err)
		ctx.Data(iris.StatusBadRequest, nil)
		return
	}

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

	users, err := db.GetUsers(context.Background())
	if err != nil {
		resp.ErrorMsg = err.Error()
		pbResponse(ctx, &resp)
		return
	}

	for _, u := range users {
		resp.Users = append(resp.Users, u.ToUserInfo())
	}
	pbResponse(ctx, &resp)
}

func getCurrentUserHandler(ctx *iris.Context) {
	req := pb.GetCurrentUserRequest{}
	resp := pb.GetCurrentUserResponse{}

	if err := proto.Unmarshal(ctx.PostBody(), &req); err != nil {
		glog.Errorf("Failed to unmarshal GetUserRequest, %v", err)
		ctx.Data(iris.StatusBadRequest, nil)
		return
	}

	curUser, ok := ctx.Session().Get(sessionUserKey).(*db.User)
	if !ok {
		resp.ErrorMsg = "Not login."
		pbResponse(ctx, &resp)
		return
	}

	resp.User = curUser.ToUserInfo()
	pbResponse(ctx, &resp)
}

func upsertUserHandler(ctx *iris.Context) {
	req := pb.UpsertUserRequest{}
	resp := pb.UpsertUserResponse{}

	if err := proto.Unmarshal(ctx.PostBody(), &req); err != nil {
		glog.Errorf("Failed to unmarshal UpsertUserRequest, %v", err)
		ctx.Data(iris.StatusBadRequest, nil)
		return
	}

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

	if curUser.Disabled {
		resp.ErrorMsg = "Disabled."
		pbResponse(ctx, &resp)
		return
	}

	newUser := db.FromUserInfo(req.User)
	if err := db.UpsertUser(context.Background(), curUser, newUser, req.IsNew); err != nil {
		glog.Errorf("Failed to upsert user, %v", err)
		resp.ErrorMsg = "Failed to save user, please reload and try again."
	}

	pbResponse(ctx, &resp)
}

func isAdmin(user *db.User) bool {
	rs := user.GetRoles()
	_, ok := rs[db.RoleAdmin]
	return ok
}
