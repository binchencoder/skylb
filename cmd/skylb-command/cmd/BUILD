package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "skylb-command",
    srcs = [
        "conf.go",
        "main.go",
        "sort.go",
    ],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_peterh_liner//:go_default_library",
        "@io_k8s_api//core/v1:go_default_library",
    ],
)
