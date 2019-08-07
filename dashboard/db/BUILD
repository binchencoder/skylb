package(default_visibility = ["//skylb:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "conf.go",
        "logs.go",
        "tx.go",
        "users.go",
    ],
    importpath = "github.com/binchencoder/skylb/dashboard/db",
    deps = [
        "//letsgo/ldap:go_default_library",
        "//skylb/dashboard/proto:go_default_library",
        "//skylb/dashboard/util:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/linuxerwang/confish:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/gopkg.in/ldap.v2:go_default_library",
        "//third-party-go/vendor/upper.io/db.v3:go_default_library",
        "//third-party-go/vendor/upper.io/db.v3/lib/sqlbuilder:go_default_library",
        "//third-party-go/vendor/upper.io/db.v3/mysql:go_default_library",
        "//third-party-go/vendor/upper.io/db.v3/postgresql:go_default_library",
    ],
)
