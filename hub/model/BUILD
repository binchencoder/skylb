package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "endpoint.go",
        "mock.go",
        "observer.go",
        "service.go",
    ],
    importpath = "github.com/binchencoder/skylb/hub/model",
    deps = [
        "//skylb-api/proto:go_default_library",
        "//skylb/hub/util:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/stretchr/testify/mock:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "endpoint_test.go",
        "observer_test.go",
        "service_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    deps = [
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
