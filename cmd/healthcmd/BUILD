package(default_visibility = ["//skylb:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "grpchealth",
    srcs = ["main.go"],
    deps = [
        "//letsgo:go_default_library",
        "//letsgo/time:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/health/grpc_health_v1:go_default_library",
    ],
)

go_binary(
    name = "svchealth",
    srcs = [
        "svchealth.go",
    ],
    deps = [
        "//letsgo:go_default_library",
        "//letsgo/service/naming:go_default_library",
        "//skylb-api/server:go_default_library",
        "//skylb/cmd/grpchealth:go_default_library",
        "//skylb/hub:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
    ],
)
