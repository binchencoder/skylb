load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
# load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test")

go_binary(
    name = "skylb",
    srcs = ["main.go"],
    deps = [
        "//rpc:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//grpc:go_default_library",
        "@com_github_binchencoder_letsgo//metrics:go_default_library",
        "@com_github_binchencoder_letsgo//runtime/pprof:go_default_library",
        "@com_github_binchencoder_skylb_api//metrics:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",        
        "@com_github_golang_glog//:go_default_library",
        "@com_github_soheilhy_cmux//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//health:go_default_library",
        "@org_golang_google_grpc//health/grpc_health_v1:go_default_library",
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
        "//rpc:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//grpc:go_default_library",
        "@com_github_binchencoder_letsgo//metrics:go_default_library",
        "@com_github_binchencoder_letsgo//runtime/pprof:go_default_library",
        "@com_github_binchencoder_skylb_api//metrics:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",        
        "@com_github_golang_glog//:go_default_library",
        "@com_github_soheilhy_cmux//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//health:go_default_library",
        "@org_golang_google_grpc//health/grpc_health_v1:go_default_library",
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
