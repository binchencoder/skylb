package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "endpoints.go",
        "hub.go",
        "int_test_common.go",
        "k8s.go",
        "key.go",
        "observer.go",
        "svcgraph.go",
    ],
    importpath = "binchencoder.com/skylb/hub",
    deps = [
        "@com_github_binchencoder_letsgo//strings:go_default_library",
        "@com_github_binchencoder_letsgo//sync:go_default_library",
        "@com_github_binchencoder_skylb_api//lameduck:go_default_library",
        "@com_github_binchencoder_skylb_api//prefix:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",
        "@com_github_binchencoder_skylb_api//util:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_prometheus_client_golang//prometheus:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@io_k8s_api//core/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/fields:go_default_library",
        "@io_k8s_client_go//kubernetes:go_default_library",
        "@io_k8s_client_go//kubernetes/typed/core/v1:go_default_library",
        "@io_k8s_client_go//rest:go_default_library",
        "@io_k8s_client_go//tools/cache:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "endpoints_test.go",
        "hub_test.go",
        "key_test.go",
        "observer_test.go",
        "svcgraph_com_test.go",
        "svcgraph_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    deps = [
        "@com_github_binchencoder_letsgo//testing/mocks/etcd:go_default_library",
        "@com_github_binchencoder_skylb_api//proto:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_smartystreets_goconvey//convey:go_default_library",
        "@io_k8s_api//core/v1:go_default_library",
    ],
)

go_test(
    name = "sg_test",
    srcs = ([
        # The order of following src files are essential for seeing useful results.
        "svcgraph_test.go",
        "svcgraph_com_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    deps = ["@com_github_golang_glog//:go_default_library"],
)

go_test(
    name = "int_test",
    srcs = ([
        # The order of following src files are essential for seeing useful results.
        "svcgraph_int_test.go",
        "svcgraph_com_test.go",
        "hub_int_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    tags = ["requires-network"],
    deps = [
        "@com_github_binchencoder_skylb_api//proto:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
)
