package(default_visibility = ["//:__subpackages__"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "static",
    deps = [
        "/dashboard/static/images",
        "/dashboard/static/scripts",
        "/dashboard/static/styles",
    ],
)
