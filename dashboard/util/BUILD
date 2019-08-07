package(default_visibility = ["//skylb:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "service.go",
    ],
    importpath = "github.com/binchencoder/skylb/dashboard/util",
    deps = ["//ease-gateway/proto/data:go_default_library"],
)
