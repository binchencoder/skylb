package(default_visibility = ["//skylb:__subpackages__"])

load("//bld_tools/bazel/rules_jingoal/protobuf:def.bzl", "genproto_go", "genproto_js")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

genproto_js(
    name = "proto_jssrc",
    srcs = [
        "api.proto",
        "data.proto",
    ],
    library = "skylb",
)

closure_js_library(
    name = "proto_js",
    srcs = [":proto_jssrc"],
    deps = [
        "//third-party-js/protobuf/jspb:jspb_js",
        "//third-party-js/protobuf/jspb/binary:binary_js",
        "@io_bazel_rules_closure//closure/library",
    ],
)

genproto_go(
    name = "proto_gosrc",
    srcs = [
        "api.proto",
        "data.proto",
    ],
    has_service = True,
)

go_library(
    name = "go_default_library",
    srcs = [":proto_gosrc"],
    importpath = "github.com/binchencoder/skylb/dashboard/proto",
    deps = [
        "//third-party-go/vendor/github.com/golang/protobuf/proto:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
    ],
)
