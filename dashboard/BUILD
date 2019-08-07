package(default_visibility = ["//skylb:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"]),
    importpath = "github.com/binchencoder/skylb/dashboard",
    deps = [
        "//letsgo/service/naming:go_default_library",
        "//letsgo/strings:go_default_library",
        "//skylb-api/lameduck:go_default_library",
        "//skylb-api/prefix:go_default_library",
        "//skylb/dashboard/db:go_default_library",
        "//skylb/dashboard/proto:go_default_library",
        "//skylb/dashboard/templates:dashboard_html",
        "//skylb/dashboard/templates:login_html",
        "//skylb/dashboard/util:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/gogo/protobuf/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/golang/protobuf/proto:go_default_library",
        "//third-party-go/vendor/github.com/kataras/iris:go_default_library",
        "//third-party-go/vendor/github.com/linuxerwang/goats-html/runtime:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
    ],
)
