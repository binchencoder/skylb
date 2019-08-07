package(default_visibility = ["//skylb:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "github.com/binchencoder/skylb/rpc",
    deps = [
        "//skylb-api/lameduck:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//skylb/hub:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/peer:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "server_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    deps = [
        "//skylb-api/proto:go_default_library",
        "//skylb/hub:go_default_library",
        "//third-party-go/vendor/github.com/stretchr/testify/mock:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/metadata:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/peer:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
    ],
)
