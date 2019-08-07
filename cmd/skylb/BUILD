load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
# load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test")

go_binary(
    name = "skylb",
    srcs = ["main.go"],
    deps = [
        "//letsgo:go_default_library",
        "//letsgo/grpc:go_default_library",
        "//letsgo/metrics:go_default_library",
        "//letsgo/runtime/pprof:go_default_library",
        "//skylb-api/metrics:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//skylb/rpc:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/soheilhy/cmux:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/health:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/health/grpc_health_v1:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "main.go",
        "main_test.go",
    ]),
    deps = [
        "//letsgo:go_default_library",
        "//letsgo/grpc:go_default_library",
        "//letsgo/metrics:go_default_library",
        "//letsgo/runtime/pprof:go_default_library",
        "//skylb-api/metrics:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//skylb/rpc:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/soheilhy/cmux:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/health:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/health/grpc_health_v1:go_default_library",
    ],
)

pkg_tar(
    name = "skylb_tar",
    srcs = [
        "start.sh",
        ":skylb",
    ],
    package_dir = "/skylb/bin",
)

# container_image(
#     name = "latest",
#     base = "@base//image",
#     tars = [
#         ":skylb_tar",
#         "//production/config/vexconfig",
#     ],
# )

# container_push(
#     name = "release",
#     format = "Docker",
#     image = ":latest",
#     registry = "",
#     repository = "skylb/skylb",
# )
