package(default_visibility = ["//skylb:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "github.com/binchencoder/skylb/cmd/webserver/svclist",
    deps = [
        "//skylb-api/prefix:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//skylb/hub:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/k8s.io/api/core/v1:go_default_library",
    ],
)

go_test(
    name = "int_test",
    srcs = glob([
        "*_int_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    tags = ["requires-network"],
    deps = [
        "//skylb-api/proto:go_default_library",
        "//skylb/hub:go_default_library",
    ],
)
