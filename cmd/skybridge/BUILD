# load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test")

go_binary(
    name = "skybridge",
    srcs = [
        "bridge.go",
        "main.go",
    ],
    deps = [
        "//letsgo:go_default_library",
        "//letsgo/strings:go_default_library",
        "//skylb-api/util:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "bridge.go",
        "main.go",
        "main_test.go",
    ]),
    deps = [
        "//letsgo:go_default_library",
        "//letsgo/strings:go_default_library",
        "//skylb-api/util:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
    ],
)

# docker_build(
#     name = "latest",
#     base = "//bld_tools/docker/ubuntu:base",
#     cmd = "/skybridge/bin/skybridge",
#     directory = "/skybridge/bin",
#     files = [
#         ":skybridge",
#     ],
#     repository = "harbor.jingoal.ltd",
#     tag = "skylb/skylbridge",
# )
