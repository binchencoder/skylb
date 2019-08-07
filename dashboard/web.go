package dashboard

import (
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/iris"
	"github.com/linuxerwang/goats-html/runtime"
)

const (
	sessionUserKey = "login_user"
)

var (
	debugMode bool
	staticDir string
	settings  runtime.TemplateSettings
)

// Init initializes the webserver.
func Init(debug bool, sDir string, etcdEps string) {
	debugMode = debug
	staticDir = sDir

	initEtcdClient(etcdEps)

	runtime.InitGoats(nil)
	settings = runtime.TemplateSettings{
		OmitDocType: false,
		DebugMode:   false,
	}

	// Login not required.
	iris.Get("/login", loginGetHandler)
	iris.Get("/logout", logoutHandler)
	iris.Post("/login", loginPostHandler)

	// Login required.
	iris.Get("/", authenticate, homeHandler)
	iris.Get("/logs", authenticate, homeHandler)
	iris.Get("/service/*path", authenticate, homeHandler)
	iris.Get("/users", authenticate, homeHandler)

	// Apis are all login required.
	api := iris.Party("_", authenticate)
	api.Post("/add-instance", addInstanceHandler)
	api.Post("/get-all-services", getAllServicesHandler)
	api.Post("/get-current-user", getCurrentUserHandler)
	api.Post("/get-logs", getLogsHandler)
	api.Post("/get-service-by-id", getServiceByIdHandler)
	api.Post("/get-users", getUsersHandler)
	api.Post("/login", loginApiHandler)
	api.Post("/toggle-lameduck", toggleLameduckHandler)
	api.Post("/upsert-user", upsertUserHandler)
}

func authenticate(ctx *iris.Context) {
	if ctx.PathString() == "/_/login" {
		ctx.Next()
		return
	}

	u := ctx.Session().Get(sessionUserKey)
	if u != nil {
		ctx.Next()
	} else {
		if strings.ToTitle(ctx.MethodString()) == "POST" {
			ctx.Text(iris.StatusUnauthorized, "Unauthorized")
		} else {
			ctx.Redirect("/login")
		}
	}
}

func handleMarkdown(ctx *iris.Context) {
	path := ctx.Param("path")
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	ctx.MustRender(path, nil)
}

func pbResponse(ctx *iris.Context, pb proto.Message) {
	b, err := proto.Marshal(pb)
	if err != nil {
		glog.Errorf("Failed to marshal pb response, %+v", err)
		ctx.Data(iris.StatusInternalServerError, nil)
		return
	}

	ctx.SetHeader("Content-Type", "application/x-protobuf")
	ctx.Data(iris.StatusOK, b)
}
