package(default_visibility = ["//:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "github.com/binchencoder/skylb/rpc",
    deps = [
        "//hub:go_default_library",
        "@com_github_binchencoder_skylb_api//lameduck:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_prometheus_client_golang//prometheus:go_default_library",
        "@org_golang_google_grpc//peer:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "server_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    deps = [
        "//hub:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",
        "@com_github_stretchr_testify//mock:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//peer:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
    ],
)
