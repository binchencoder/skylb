package dashboard

import (
	"github.com/golang/glog"
	"github.com/kataras/iris"

	"binchencoder.com/skylb/dashboard/db"
)

func logoutHandler(ctx *iris.Context) {
	_, ok := ctx.Session().Get(sessionUserKey).(*db.User)
	if !ok {
		glog.Errorf("Try to logout a user who didn't login.")
	}

	ctx.Session().Delete(sessionUserKey)

	ctx.Redirect("/login")
}
