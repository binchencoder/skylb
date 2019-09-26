# load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test")

go_binary(
    name = "skylbweb",
    srcs = ["main.go"],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//metrics:go_default_library",
        "@com_github_binchencoder_letsgo//strings:go_default_library",
        "@com_github_binchencoder_skylb_api//lameduck:go_default_library",
        "@com_github_binchencoder_skylb_api//prefix:go_default_library",
        "/cmd/grpchealth:go_default_library",
        "/cmd/webserver/svclist:go_default_library",
        "/hub:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_prometheus_client_golang//prometheus:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "main.go",  # Must include main.go so as to check flags.
        "main_test.go",
    ]),
    # The deps shall be the same as that of the go_binary.
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//metrics:go_default_library",
        "@com_github_binchencoder_letsgo//strings:go_default_library",
        "@com_github_binchencoder_skylb_api//lameduck:go_default_library",
        "@com_github_binchencoder_skylb_api//prefix:go_default_library",
        "/cmd/grpchealth:go_default_library",
        "/cmd/webserver/svclist:go_default_library",
        "/hub:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_prometheus_client_golang//prometheus:go_default_library",
    ],
)

# docker_build(
#     name = "latest",
#     base = "//bld_tools/docker/ubuntu:base",
#     cmd = "/skylb/start-webserver.sh",
#     directory = "/skylb",
#     files = [
#         "start-webserver.sh",
#         ":skylbweb",
#     ],
#     ports = [
#         "8090",
#     ],
#     repository = "harbor.jingoal.ltd",
#     tag = "skylb/webserver",
#     workdir = "/skylb",
# )
