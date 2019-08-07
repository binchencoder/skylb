package(default_visibility = ["//skylb:__subpackages__"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_binary", "closure_js_library")

closure_js_library(
    name = "main_js",
    srcs = [
        "extern.js",
        "homeview.js",
        "logs.js",
        "main.js",
        "serviceview.js",
        "usermgr.js",
        "view.js",
        "viewmgr.js",
        "xhrmgr.js",
    ],
    suppress = [
        "CR_NOT_PROVIDED",
        "superfluousSuppress",
    ],
    deps = [
        "//skylb/dashboard/proto:proto_js",
        "//skylb/dashboard/templates:home_html_js",
        "//skylb/dashboard/templates:loginbox_html_js",
        "//skylb/dashboard/templates:logs_html_js",
        "//skylb/dashboard/templates:service_html_js",
        "//skylb/dashboard/templates:usermgr_html_js",
        "@io_bazel_rules_closure//closure/library",
        "@io_bazel_rules_closure//closure/library/ui:all_js",
    ],
)

closure_js_binary(
    name = "skylb-0.1.0.min",
    entry_points = ["goog:skylb.main"],
    suppress_on_all_sources_in_transitive_closure = [
        "CR_NOT_PROVIDED",
        "JSC_EXTRA_REQUIRE_WARNING",
        "JSC_IMPLICITLY_NULLABLE_JSDOC",
        "JSC_INEXISTENT_PROPERTY",
        "JSC_MISSING_CONST_PROPERTY",
        "JSC_MISSING_PROVIDE_ERROR",
        "JSC_NULLABLE_RETURN_WITH_NAME",
        "JSC_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL",
        "JSC_REQUIRES_NOT_SORTED",
        "JSC_STRICT_INEXISTENT_PROPERTY",
        "JSC_TYPE_MISMATCH",
        "JSC_UNDEFINED_NAME",
        "JSC_UNKNOWN_EXPR_TYPE",
        "JSC_UNRECOGNIZED_TYPE_ERROR",
        "JSC_UNUSED_LOCAL_ASSIGNMENT",
        "JSC_UNUSED_PRIVATE_PROPERTY",
        "JSC_VAR_MULTIPLY_DECLARED_ERROR",
        "deprecated",
        "lintChecks",
        "strictMissingRequire",
        "superfluousSuppress",
        "unusedLocalVariables",
    ],
    deps = [
        ":main_js",
    ],
)

closure_js_binary(
    name = "skylb-0.1.0.min-dev",
    compilation_level = "WHITESPACE_ONLY",
    debug = True,
    entry_points = ["goog:skylb.main"],
    formatting = "PRETTY_PRINT",
    deps = [
        ":main_js",
    ],
)

pkg_tar(
    name = "skylb",
    srcs = [":skylb-0.1.0.min"],
    package_dir = "/skylb/dashboard/static/scripts",
)
