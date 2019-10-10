package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "binchencoder.com/skylb/cmd/webserver/svclist",
    deps = [
        "//hub:go_default_library",
        "@com_github_binchencoder_skylb_api//prefix:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",        
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@io_k8s_api//core/v1:go_default_library",
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
        "//hub:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",
    ],
)
