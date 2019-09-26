load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
# load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test")

go_binary(
    name = "skylb-dashboard",
    srcs = ["main.go"],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "/dashboard:go_default_library",
        "/dashboard/db:go_default_library",
        "@com_github_kataras_iris//:go_default_library",
    ],
)

go_binary(
    name = "etcd-loader",
    srcs = [
        "data.go",
        "etcdloader.go",
    ],
    deps = [
        "@com_github_binchencoder_letsgo//strings:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)

pkg_tar(
    name = "skylb-dashboard-tar",
    srcs = [
        ":skylb-dashboard",
    ],
    package_dir = "/skylb/dashboard/bin",
)

pkg_tar(
    name = "etcd-loader-tar",
    srcs = [
        ":etcd-loader",
    ],
    package_dir = "/skylb/dashboard/bin",
)

pkg_tar(
    name = "sql",
    srcs = [
        "1000_schema.up.sql",
        "2000_data.up.sql",
    ],
    package_dir = "/dataloader/sql",
)

pkg_tar(
    name = "docs-sql",
    srcs = [
        "1000_schema.up.sql",
        "2000_data.up.sql",
    ],
    package_dir = "/skylb/dashboard/docs/sql",
)

pkg_tar(
    name = "supervisord-tar",
    srcs = ["supervisord.conf"],
    package_dir = "/dataloader",
)

pkg_tar(
    name = "shell-scripts",
    srcs = [
        "start.sh",
        "stop.sh",
    ],
    package_dir = "/skylb/dashboard",
)

pkg_tar(
    name = "release-tar",
    deps = [
        ":shell-scripts",
        ":skylb-dashboard-tar",
        "/dashboard/scripts/skylb",
        "/dashboard/static",
        "/dashboard/styles",
        "/docs/dashboard",
    ],
)

pkg_tar(
    name = "dev-data-tar",
    deps = [
        ":etcd-loader-tar",
        ":sql",
    ],
)

# docker_build(
#     name = "latest",
#     base = "//docker-compose/dataloader/cmd/dataloader:latest",
#     ports = [
#         "8050",
#     ],
#     repository = "harbor.jingoal.ltd",
#     tag = "skylb/dashboard",
#     tars = [
#         ":dev-data-tar",
#         ":release-tar",
#         ":supervisord-tar",
#     ],
# )
