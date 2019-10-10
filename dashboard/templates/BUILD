package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//bld_tools/bazel/rules_jingoal/goats-html:def.bzl", "gengoats_go", "gengoats_closure")

gengoats_go(
    name = "dashboard_html_gosrc",
    src = "dashboard.html",
)

go_library(
    name = "dashboard_html",
    srcs = [":dashboard_html_gosrc"],
    importpath = "binchencoder.com/skylb/dashboard/templates/dashboard_html",
    deps = [
        "/dashboard/proto:go_default_library",
        "@com_github_linuxerwang_goats_html//runtime:go_default_library",
    ],
)

gengoats_closure(
    name = "home_html_jssrc",
    src = "home.html",
    pkg_prefix = "skylb",
)

closure_js_library(
    name = "home_html_js",
    srcs = [
        ":home_html_jssrc",
    ],
    deps = [
        "@com_github_linuxerwang_goats_html//client/goats/runtime",
        "@com_github_linuxerwang_goats_html//client/goats/runtime/filters",
        "@io_bazel_rules_closure//closure/library",
    ],
)

gengoats_go(
    name = "loginbox_html_gosrc",
    src = "loginbox.html",
)

go_library(
    name = "loginbox_html",
    srcs = [":loginbox_html_gosrc"],
    importpath = "binchencoder.com/skylb/dashboard/templates/loginbox_html",
    deps = [
        "@com_github_linuxerwang_goats_html//runtime:go_default_library",
    ],
)

gengoats_go(
    name = "login_html_gosrc",
    src = "login.html",
    deps = [
        ":loginbox_html_gosrc",
    ],
)

go_library(
    name = "login_html",
    srcs = [":login_html_gosrc"],
    importpath = "binchencoder.com/skylb/dashboard/templates/login_html",
    deps = [
        ":loginbox_html",
        "@com_github_linuxerwang_goats_html//runtime:go_default_library",
    ],
)

gengoats_closure(
    name = "loginbox_html_jssrc",
    src = "loginbox.html",
    pkg_prefix = "skylb",
)

closure_js_library(
    name = "loginbox_html_js",
    srcs = [
        ":loginbox_html_jssrc",
    ],
    deps = [
        "@com_github_linuxerwang_goats_html//client/goats/runtime",
        "@com_github_linuxerwang_goats_html//client/goats/runtime/filters",
        "@io_bazel_rules_closure//closure/library",
    ],
)

gengoats_closure(
    name = "logs_html_jssrc",
    src = "logs.html",
    pkg_prefix = "skylb",
)

closure_js_library(
    name = "logs_html_js",
    srcs = [
        ":logs_html_jssrc",
    ],
    deps = [
        "@com_github_linuxerwang_goats_html//client/goats/runtime",
        "@com_github_linuxerwang_goats_html//client/goats/runtime/filters",
        "@io_bazel_rules_closure//closure/library",
    ],
)

gengoats_closure(
    name = "usermgr_html_jssrc",
    src = "usermgr.html",
    pkg_prefix = "skylb",
)

closure_js_library(
    name = "usermgr_html_js",
    srcs = [
        ":usermgr_html_jssrc",
    ],
    deps = [
        "@com_github_linuxerwang_goats_html//client/goats/runtime",
        "@com_github_linuxerwang_goats_html//client/goats/runtime/filters",
        "@io_bazel_rules_closure//closure/library",
    ],
)

gengoats_closure(
    name = "service_html_jssrc",
    src = "service.html",
    pkg_prefix = "skylb",
)

closure_js_library(
    name = "service_html_js",
    srcs = [
        ":service_html_jssrc",
    ],
    deps = [
        "@com_github_linuxerwang_goats_html//client/goats/runtime",
        "@com_github_linuxerwang_goats_html//client/goats/runtime/filters",
        "@io_bazel_rules_closure//closure/library",
    ],
)
