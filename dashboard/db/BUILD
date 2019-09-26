package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "conf.go",
        "logs.go",
        "tx.go",
        "users.go",
    ],
    importpath = "github.com/binchencoder/skylb/dashboard/db",
    deps = [
        "@com_github_binchencoder_letsgo//ldap:go_default_library",
        "/dashboard/proto:go_default_library",
        "/dashboard/util:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_linuxerwang_confish//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@in_gopkg_ldap_v2//:go_default_library",
        "@io_upper_db_v3//:go_default_library",
        "@io_upper_db_v3//lib/sqlbuilder:go_default_library",
        "@io_upper_db_v3//mysql:go_default_library",
        "@io_upper_db_v3//postgresql:go_default_library",
    ],
)
