package dashboard

import (
	"github.com/golang/glog"
	"github.com/kataras/iris"

	db "binchencoder.com/skylb/dashboard/db"
	tpl "binchencoder.com/skylb/dashboard/templates/dashboard_html"
)

func homeHandler(ctx *iris.Context) {
	ctx.SetHeader("Content-Type", "text/html")

	curUser, ok := ctx.Session().Get(sessionUserKey).(*db.User)
	if !ok {
		ctx.NotFound()
		return
	}

	t := tpl.NewMainPageTemplate(ctx.Response.BodyWriter(), &settings)
	args := &tpl.MainPageTemplateArgs{
		User: curUser.ToUserInfo(),
	}
	err := t.Render(args)
	if err != nil {
		glog.Errorln("Error to render template:", err)
	}
}
