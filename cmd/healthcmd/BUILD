package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "grpchealth",
    srcs = ["main.go"],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//time:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//health/grpc_health_v1:go_default_library",
    ],
)

go_binary(
    name = "svchealth",
    srcs = [
        "svchealth.go",
    ],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//service/naming:go_default_library",
        "@com_github_binchencoder_skylb_api//server:go_default_library",
        "/cmd/grpchealth:go_default_library",
        "/hub:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
    ],
)
