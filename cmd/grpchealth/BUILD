package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["healthlib.go"],
    ),
    importpath = "binchencoder.com/skylb/cmd/grpchealth",
    deps = [
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
        "@com_github_binchencoder_letsgo//service/naming:go_default_library",
        "@com_github_binchencoder_skylb_api//client:go_default_library",
        "@com_github_binchencoder_skylb_api//internal/health:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",
        "//cmd/webserver/svclist:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_prometheus_client_golang//prometheus:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//health/grpc_health_v1:go_default_library",
    ],
)
