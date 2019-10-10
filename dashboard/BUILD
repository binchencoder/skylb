package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"]),
    importpath = "binchencoder.com/skylb/dashboard",
    deps = [
        "@com_github_binchencoder_letsgo//service/naming:go_default_library",
        "@com_github_binchencoder_letsgo//strings:go_default_library",
        "@com_github_binchencoder_skylb_api//lameduck:go_default_library",
        "@com_github_binchencoder_skylb_api//prefix:go_default_library",
        "/dashboard/db:go_default_library",
        "/dashboard/proto:go_default_library",
        "/dashboard/templates:dashboard_html",
        "/dashboard/templates:login_html",
        "/dashboard/util:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_kataras_iris//:go_default_library",
        "@com_github_linuxerwang_goats_html//runtime:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
    ],
)
