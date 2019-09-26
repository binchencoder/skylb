package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "key.go",
    ],
    importpath = "github.com/binchencoder/skylb/hub/util",
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "key_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    deps = ["@com_github_smartystreets_goconvey//convey:go_default_library"],
)
