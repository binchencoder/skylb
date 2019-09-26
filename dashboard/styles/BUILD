package(default_visibility = ["//:__subpackages__"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_css_binary", "closure_css_library")

closure_css_library(
    name = "skylb_css",
    srcs = [
        "skylb.css",
    ],
)

closure_css_binary(
    name = "skylb-0.1.0.min",
    renaming = 0,
    deps = [
        ":skylb_css",
        "@io_bazel_rules_closure//closure/library:css",
    ],
)

pkg_tar(
    name = "styles",
    srcs = [
        ":skylb-0.1.0.min",
    ],
    package_dir = "/skylb/dashboard/static/styles",
)
