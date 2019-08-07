package(default_visibility = ["//skylb:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["healthlib.go"],
    ),
    importpath = "github.com/binchencoder/skylb/cmd/grpchealth",
    deps = [
        "//letsgo/service/naming:go_default_library",
        "//skylb-api/client:go_default_library",
        "//skylb-api/internal/health:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//skylb/cmd/webserver/svclist:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/health/grpc_health_v1:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
    ],
)
