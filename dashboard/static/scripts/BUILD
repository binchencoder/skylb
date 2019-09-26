package(default_visibility = ["//:__subpackages__"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "scripts",
    srcs = glob(
        ["*"],
        exclude = ["BUILD"],
    ),
    package_dir = "/skylb/dashboard/static/scripts",
)
