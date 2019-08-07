package(default_visibility = ["//skylb:__subpackages__"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "certs",
    srcs = [
        "cert.pem",
        "key.pem",
    ],
    package_dir = "/skylb/dashboard/certs",
)

pkg_tar(
    name = "conf",
    srcs = ["database.conf"],
    package_dir = "/skylb/dashboard/conf",
)

pkg_tar(
    name = "dashboard",
    deps = [
        ":certs",
        ":conf",
    ],
)
