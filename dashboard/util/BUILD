package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "service.go",
    ],
    importpath = "github.com/binchencoder/skylb/dashboard/util",
    deps = ["@com_github_binchencoder_gateway_proto//data:go_default_library"],
)
