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
    importpath = "github.com/binchencoder/skylb/hub",
    deps = [
        "//letsgo/strings:go_default_library",
        "//letsgo/sync:go_default_library",
        "//skylb-api/lameduck:go_default_library",
        "//skylb-api/prefix:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//skylb-api/util:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/k8s.io/api/core/v1:go_default_library",
        "//third-party-go/vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//third-party-go/vendor/k8s.io/apimachinery/pkg/fields:go_default_library",
        "//third-party-go/vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//third-party-go/vendor/k8s.io/client-go/kubernetes/typed/core/v1:go_default_library",
        "//third-party-go/vendor/k8s.io/client-go/rest:go_default_library",
        "//third-party-go/vendor/k8s.io/client-go/tools/cache:go_default_library",
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
        "//letsgo/testing/mocks/etcd:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//third-party-go/vendor/k8s.io/api/core/v1:go_default_library",
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
    deps = ["//third-party-go/vendor/github.com/golang/glog:go_default_library"],
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
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
    ],
)
