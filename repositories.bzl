load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    # go_repository(
    #     name = "com_github_binchencoder_letsgo",
    #     importpath = "binchencoder.com/letsgo",
    #     urls = [
    #         "https://codeload.github.com/binchencoder/letsgo/tar.gz/3a34eef5d1546b1be444e2e053d425e80afe100e",
    #     ],
    #     strip_prefix = "letsgo-3a34eef5d1546b1be444e2e053d425e80afe100e",
    #     type = "tar.gz",
    #     # gazelle args: -go_prefix binchencoder.com/letsgo
    # )
    # go_repository(
    #     name = "com_github_binchencoder_skylb_api",
    #     importpath = "binchencoder.com/skylb-api",
    #     urls = [
    #         "https://codeload.github.com/binchencoder/skylb-api/tar.gz/50ba73497500a4239c089a8fb0772d0fd23cc1a1",
    #     ],
    #     strip_prefix = "skylb-api-50ba73497500a4239c089a8fb0772d0fd23cc1a1",
    #     type = "tar.gz",
    #     # gazelle args: -go_prefix binchencoder.com/skylb-api
    # )
    # go_repository(
    #     name = "com_github_binchencoder_gateway_proto",
    #     importpath = "binchencoder.com/gateway-proto",
    #     urls = [
    #         "https://codeload.github.com/binchencoder/gateway-proto/tar.gz/1a34b562e1a85e05485b51fa49833d5c3e0bca89",
    #     ],
    #     strip_prefix = "gateway-proto-1a34b562e1a85e05485b51fa49833d5c3e0bca89",
    #     type = "tar.gz",
    #     # gazelle args: -go_prefix binchencoder.com/gateway-proto
    # )

    go_repository(
        name = "com_github_beorn7_perks",
        importpath = "github.com/beorn7/perks",
        urls = ["https://codeload.github.com/beorn7/perks/tar.gz/37c8de3658fcb183f997c4e13e8337516ab753e6"],
        strip_prefix = "perks-37c8de3658fcb183f997c4e13e8337516ab753e6",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_cenkalti_backoff",
        importpath = "github.com/cenkalti/backoff",
        urls = ["https://github.com/cenkalti/backoff/archive/v2.2.1.tar.gz"],
        strip_prefix = "backoff-2.2.1",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_coreos_etcd",
        importpath = "github.com/coreos/etcd",
        urls = ["https://codeload.github.com/etcd-io/etcd/tar.gz/98d308426819d892e149fe45f6fd542464cb1f9d"],
        strip_prefix = "etcd-98d308426819d892e149fe45f6fd542464cb1f9d",
        type = "tar.gz",
        build_file_generation = "on",
    )
    go_repository(
        name = "com_github_coreos_go_semver",
        importpath = "github.com/coreos/go-semver",
        commit = "e214231b295a8ea9479f11b70b35d5acf3556d9b",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        commit = "d8f796af33cc11cb798c1aaeb27a4ebc5099927d",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
        version = "v0.0.0-20160126235308-23def4e6c14b",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        commit = "c2e93f3ae59f2904160ceaab466009f965df46d6",
        # gazelle args: -go_prefix github.com/google/uuid
    )
    go_repository(
        name = "com_github_google_gofuzz",
        importpath = "github.com/google/gofuzz",
        commit = "b906efc57a556621b61db18d73df8c109dfa3613",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        commit = "b1c9c4891a6525d98001fea424c8926c6d77bb56",
    )
    go_repository(
        name = "com_github_googleapis_gnostic",
        importpath = "github.com/googleapis/gnostic",
        urls = [
            "https://codeload.github.com/googleapis/gnostic/tar.gz/b0a17e38ce1aad0c792ef9efd1810364be151db4",
        ],
        strip_prefix = "gnostic-b0a17e38ce1aad0c792ef9efd1810364be151db4",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_gogo_protobuf",
        importpath = "github.com/gogo/protobuf",
        urls = [
            "https://codeload.github.com/gogo/protobuf/tar.gz/8a5ed79f688836cf007ca23aefe0299791e7bea5",
        ],
        strip_prefix = "protobuf-8a5ed79f688836cf007ca23aefe0299791e7bea5",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        commit = "a1ebd699b1950beb1da0752cbb5559662018f798",
    )
    go_repository(
        name = "com_github_go_kit_kit",
        importpath = "github.com/go-kit/kit",
        urls = ["https://codeload.github.com/go-kit/kit/tar.gz/dc489b75b9cdbf29c739534c2aa777cabb034954"],
        strip_prefix = "kit-dc489b75b9cdbf29c739534c2aa777cabb034954",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_opentracing",
        importpath = "github.com/grpc-ecosystem/grpc-opentracing",
        urls = [
            "https://codeload.github.com/grpc-ecosystem/grpc-opentracing/tar.gz/8e809c8a86450a29b90dcc9efbf062d0fe6d9746",
        ],
        strip_prefix = "grpc-opentracing-8e809c8a86450a29b90dcc9efbf062d0fe6d9746",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_middleware",
        importpath = "github.com/grpc-ecosystem/go-grpc-middleware",
        urls = [
            "https://codeload.github.com/grpc-ecosystem/go-grpc-middleware/tar.gz/e0797f438f94f4d032395b8f71aae0e73d6efa08",
        ],
        strip_prefix = "go-grpc-middleware-e0797f438f94f4d032395b8f71aae0e73d6efa08",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_hashicorp_golang_lru",
        importpath = "github.com/hashicorp/golang-lru",
        commit = "7f827b33c0f158ec5dfbba01bb0b14a4541fd81d",
    )
    go_repository(
        name = "com_github_json_iterator_go",
        importpath = "github.com/json-iterator/go",
        commit = "819acad769e54806c920726ac93537ba4e2c22ad",
    )
    go_repository(
        name = "com_github_jtolds_gls",
        importpath = "github.com/jtolds/gls",
        commit = "b4936e06046bbecbb94cae9c18127ebe510a2cb9",
    )
    go_repository(
        name = "com_github_kataras_iris",
        importpath = "github.com/kataras/iris",
        urls = ["https://codeload.github.com/kataras/iris/tar.gz/df882273e21952a316236174123fc09096b49aad"],
        strip_prefix = "iris-df882273e21952a316236174123fc09096b49aad",
        type = "tar.gz",
        build_file_generation = "on",
    )
    go_repository(
        name = "com_github_linuxerwang_goats_html",
        importpath = "github.com/linuxerwang/goats-html",
        commit = "cdff773a61b4faf647611ea9d73f04848c7fe096",
    )
    go_repository(
        name = "com_github_linuxerwang_confish",
        importpath = "github.com/linuxerwang/confish",
        commit = "e1f17b4f6bb632f8e5d9d73242917c1d4c723710",
    )
    go_repository(
        name = "com_github_matttproud_golang_protobuf_extensions",
        importpath = "github.com/matttproud/golang_protobuf_extensions",
        commit = "c182affec369e30f25d3eb8cd8a478dee585ae7d",
    )
    go_repository(
        name = "com_github_mattn_go_runewidth",
        importpath = "github.com/mattn/go-runewidth",
        commit = "703b5e6b11ae25aeb2af9ebb5d5fdf8fa2575211",
    )
    go_repository(
        name = "com_github_modern_go_reflect2",
        importpath = "github.com/modern-go/reflect2",
        commit = "94122c33edd36123c84d5368cfb2b69df93a0ec8",
    )
    go_repository(
        name = "com_github_modern_go_concurrent",
        importpath = "github.com/modern-go/concurrent",
        commit = "bacd9c7ef1dd9b15be4a9909b8ac7a4e313eec94",
    )
    go_repository(
        name = "com_github_opentracing_opentracing_go",
        importpath = "github.com/opentracing/opentracing-go",
        urls = [
            "https://codeload.github.com/opentracing/opentracing-go/tar.gz/135aa78c6f95b4a199daf2f0470d231136cbbd0c",
        ],
        strip_prefix = "opentracing-go-135aa78c6f95b4a199daf2f0470d231136cbbd0c",
        type = "tar.gz",
        # gazelle args: -go_prefix github.com/opentracing/opentracing-go
    )
    go_repository(
        name = "com_github_pborman_uuid",
        importpath = "github.com/pborman/uuid",
        commit = "8b1b92947f46224e3b97bb1a3a5b0382be00d31e",
        # gazelle args: -go_prefix github.com/pborman/uuid
    )
    go_repository(
        name = "com_github_peterh_liner",
        importpath = "github.com/peterh/liner",
        commit = "6f820f8f90ce9482ffbd40bb15f9ea9932f4942d",
        # gazelle args: -go_prefix github.com/peterh/liner
    )
    go_repository(
        name = "com_github_pkg_errors",
        importpath = "github.com/pkg/errors",
        commit = "27936f6d90f9c8e1145f11ed52ffffbfdb9e0af7",
    )
    go_repository(
        name = "com_github_prometheus_client_golang",
        importpath = "github.com/prometheus/client_golang",
        urls = [
            "https://codeload.github.com/prometheus/client_golang/tar.gz/b12dd9c58c3d7ce96f9e1ede31d02f6df3d50c61",
        ],
        strip_prefix = "client_golang-b12dd9c58c3d7ce96f9e1ede31d02f6df3d50c61",
        type = "tar.gz",
        # gazelle args: -go_prefix github.com/prometheus/client_golang
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        importpath = "github.com/prometheus/client_model",
        urls = [
            "https://codeload.github.com/prometheus/client_model/tar.gz/fd36f4220a901265f90734c3183c5f0c91daa0b8",
        ],
        strip_prefix = "client_model-fd36f4220a901265f90734c3183c5f0c91daa0b8",
        type = "tar.gz",
        # gazelle args: -go_prefix github.com/prometheus/client_model
    )
    go_repository(
        name = "com_github_prometheus_common",
        importpath = "github.com/prometheus/common",
        urls = [
            "https://codeload.github.com/prometheus/common/tar.gz/637d7c34db122e2d1a25d061423098663758d2d3",
        ],
        strip_prefix = "common-637d7c34db122e2d1a25d061423098663758d2d3",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_prometheus_procfs",
        importpath = "github.com/prometheus/procfs",
        urls = [
            "https://codeload.github.com/prometheus/procfs/tar.gz/6df11039f8de6804bb01c0ebd52cde9c26091e1c",
        ],
        strip_prefix = "procfs-6df11039f8de6804bb01c0ebd52cde9c26091e1c",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_soheilhy_cmux",
        importpath = "github.com/soheilhy/cmux",
        commit = "8a8ea3c53959009183d7914522833c1ed8835020",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        commit = "221dbe5ed46703ee255b1da0dec05086f5035f62",
    )
    go_repository(
        name = "com_github_smartystreets_goconvey",
        importpath = "github.com/smartystreets/goconvey",
        urls = ["https://github.com/smartystreets/goconvey/archive/1.6.3.tar.gz"],
        strip_prefix = "goconvey-1.6.3",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_smartystreets_assertions",
        importpath = "github.com/smartystreets/assertions",
        urls = ["https://github.com/smartystreets/assertions/archive/v1.0.1.tar.gz"],
        strip_prefix = "assertions-1.0.1",
        type = "tar.gz",
    )
    go_repository(
        name = "com_github_uber_jaeger_client_go",
        importpath = "github.com/uber/jaeger-client-go",
        urls = [
            "https://codeload.github.com/jaegertracing/jaeger-client-go/tar.gz/d8999ab8c9e71b2d71022f26f21bf39a3c428301",
        ],
        strip_prefix = "jaeger-client-go-d8999ab8c9e71b2d71022f26f21bf39a3c428301",
        type = "tar.gz",
        # gazelle args: -go_prefix github.com/uber/jaeger-client-go
    )
    go_repository(
        name = "com_github_uber_jaeger_lib",
        importpath = "github.com/uber/jaeger-lib",
        urls = [
            "https://codeload.github.com/jaegertracing/jaeger-lib/tar.gz/ec4562394c7d7c18dc238aad0fc921a4325a8b0a",
        ],
        strip_prefix = "jaeger-lib-ec4562394c7d7c18dc238aad0fc921a4325a8b0a",
        type = "tar.gz",
        # gazelle args: -go-prefix github.com/uber/jaeger-lib
    )
    go_repository(
        name = "com_github_vividcortex_gohistogram",
        importpath = "github.com/VividCortex/gohistogram",
        commit = "51564d9861991fb0ad0f531c99ef602d0f9866e6",
    )

    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        urls = [
            "https://codeload.github.com/grpc/grpc-go/tar.gz/df014850f6dee74ba2fc94874043a9f3f75fbfd8",
        ],
        strip_prefix = "grpc-go-df014850f6dee74ba2fc94874043a9f3f75fbfd8", # v1.17.0, latest as of 2019-01-15
        type = "tar.gz",
        # gazelle args: -go_prefix google.golang.org/grpc -proto disable
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        urls = [
            "https://codeload.github.com/golang/crypto/tar.gz/34f69633bfdcf9db92f698f8487115767eebef81",
        ],
        strip_prefix = "crypto-34f69633bfdcf9db92f698f8487115767eebef81",
        type = "tar.gz",
    )
    go_repository(
        name = "org_golang_x_lint",
        importpath = "golang.org/x/lint",
        urls = [
            "https://codeload.github.com/golang/lint/tar.gz/16217165b5de779cb6a5e4fc81fa9c1166fda457",
        ],
        strip_prefix = "lint-16217165b5de779cb6a5e4fc81fa9c1166fda457",
        type = "tar.gz",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        urls = [
            "https://codeload.github.com/golang/time/tar.gz/c4c64cad1fd0a1a8dab2523e04e61d35308e131e",
        ],
        strip_prefix = "time-c4c64cad1fd0a1a8dab2523e04e61d35308e131e",
        type = "tar.gz",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        urls = [
            "https://codeload.github.com/golang/oauth2/tar.gz/0f29369cfe4552d0e4bcddc57cc75f4d7e672a33",
        ],
        strip_prefix = "oauth2-0f29369cfe4552d0e4bcddc57cc75f4d7e672a33",
        type = "tar.gz",
        # gazelle args: -go_prefix golang.org/x/oauth2
    )
    go_repository(
        name = "in_gopkg_ldap_v2",
        importpath = "gopkg.in/ldap.v2",
        urls = [
            "https://codeload.github.com/go-ldap/ldap/tar.gz/bb7a9ca6e4fbc2129e3db588a34bc970ffe811a9",
        ],
        strip_prefix = "ldap-bb7a9ca6e4fbc2129e3db588a34bc970ffe811a9",
        type = "tar.gz",
        # gazelle args: -go_prefix gopkg.in/ldap.v2
    )
    go_repository(
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        replace = "github.com/go-yaml/yaml",
        sum = "h1:eZqMvILvSB6AhTa+FGXHupLRXfU8SFxBP4IW1wetpT4=",
        version = "v2.0.0-20170812160011-eb3733d160e7",
        # gazelle args: -go-prefix gopkg.in/yaml.v2
    )
    go_repository(
        name = "in_gopkg_inf_v0",
        importpath = "gopkg.in/inf.v0",
        urls = [
            "https://codeload.github.com/go-inf/inf/tar.gz/8237a9a5367b2a82f922b38d4b3676293e031763",
        ],
        strip_prefix = "inf-8237a9a5367b2a82f922b38d4b3676293e031763",
        type = "tar.gz",
    )
    go_repository(
        name = "io_k8s_apimachinery",
        importpath = "k8s.io/apimachinery",
        urls = ["https://codeload.github.com/kubernetes/apimachinery/tar.gz/62598f38f24eabad89ddd52347282202797a6de9"],
        strip_prefix = "apimachinery-62598f38f24eabad89ddd52347282202797a6de9",
        type = "tar.gz",
        build_file_proto_mode = "disable",
    )
    go_repository(
        name = "io_k8s_client_go",
        importpath = "k8s.io/client-go",
        urls = ["https://codeload.github.com/kubernetes/client-go/tar.gz/07054768d98de723f5da7fb60647eda1c0471a76"],
        strip_prefix = "client-go-07054768d98de723f5da7fb60647eda1c0471a76",
        type = "tar.gz",
    )
    go_repository(
        name = "io_k8s_api",
        importpath = "k8s.io/api",
        urls = ["https://codeload.github.com/kubernetes/api/tar.gz/d58b53da08f5430bb0f4e1154a73314e82b5b3aa"],
        strip_prefix = "api-d58b53da08f5430bb0f4e1154a73314e82b5b3aa",
        type = "tar.gz",
        build_file_proto_mode = "disable",
    )
    go_repository(
        name = "io_k8s_klog",
        importpath = "k8s.io/klog",
        urls = ["https://codeload.github.com/kubernetes/klog/tar.gz/e531227889390a39d9533dde61f590fe9f4b0035"],
        strip_prefix = "klog-e531227889390a39d9533dde61f590fe9f4b0035",
        type = "tar.gz",
    )
    go_repository(
        name = "io_k8s_utils",
        importpath = "k8s.io/utils",
        urls = ["https://codeload.github.com/kubernetes/utils/tar.gz/69764acb6e8e900b7c05296c5d3c9c19545475f9"],
        strip_prefix = "utils-69764acb6e8e900b7c05296c5d3c9c19545475f9",
        type = "tar.gz",
    )
    go_repository(
        name = "io_k8s_sigs_yaml",
        importpath = "sigs.k8s.io/yaml",
        urls = ["https://codeload.github.com/kubernetes-sigs/yaml/tar.gz/fd68e9863619f6ec2fdd8625fe1f02e7c877e480"],
        strip_prefix = "yaml-fd68e9863619f6ec2fdd8625fe1f02e7c877e480",
        type = "tar.gz",
    )
    go_repository(
        name = "io_upper_db_v3",
        importpath = "upper.io/db.v3",
        urls = [
            "https://codeload.github.com/upper/db/tar.gz/ff77bee152d24abc0668e7c6f145b329f2952657",
        ],
        strip_prefix = "db-ff77bee152d24abc0668e7c6f145b329f2952657",
        type = "tar.gz",
        # gazelle args: -go_prefix upper.io/db.v3
    )