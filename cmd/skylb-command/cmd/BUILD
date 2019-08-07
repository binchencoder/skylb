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
        "//letsgo:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/peterh/liner:go_default_library",
        "//third-party-go/vendor/k8s.io/api/core/v1:go_default_library",
    ],
)
