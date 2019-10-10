package dashboard

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/iris"
	"golang.org/x/net/context"

	"binchencoder.com/skylb/dashboard/db"
	pb "binchencoder.com/skylb/dashboard/proto"
	tpl "binchencoder.com/skylb/dashboard/templates/login_html"
)

func loginGetHandler(ctx *iris.Context) {
	renderLoginPage(ctx)
}

func loginPostHandler(ctx *iris.Context) {
	loginName := ctx.FormValueString("loginname")

	if debugMode {
		setDebugUser(ctx, loginName)
		return
	}

	user, err := db.Authenticate(loginName, ctx.FormValueString("password"))
	if err != nil {
		renderLoginPage(ctx)
		return
	}

	ctx.Session().Set(sessionUserKey, user)

	ctx.Redirect("/")
}

func loginApiHandler(ctx *iris.Context) {
	req := pb.LoginRequest{}
	resp := pb.LoginResponse{}

	if err := proto.Unmarshal(ctx.PostBody(), &req); err != nil {
		glog.Error("Failed to unmarshal LoginRequest", err)
		ctx.Data(iris.StatusBadRequest, nil)
		return
	}

	if debugMode {
		setDebugUser(ctx, req.LoginName)
		return
	}

	user, err := db.Authenticate(req.LoginName, req.Password)
	if err != nil {
		resp.ErrorMsg = "Login name and password not match."
		pbResponse(ctx, &resp)
		return
	}

	ctx.Session().Set(sessionUserKey, user)
	pbResponse(ctx, &resp)
}

func renderLoginPage(ctx *iris.Context) {
	ctx.SetHeader("Content-Type", "text/html")

	t := tpl.NewLoginPageTemplate(ctx.Response.BodyWriter(), &settings)
	args := &tpl.LoginPageTemplateArgs{}
	err := t.Render(args)
	if err != nil {
		glog.Errorln("Error to render template:", err)
	}
}

func setDebugUser(ctx *iris.Context, loginName string) {
	user := db.User{LoginName: loginName}
	db.MergeUserRoles(context.Background(), &user)
	ctx.Session().Set(sessionUserKey, &user)
	ctx.Redirect("/")
}
