package hub

import (
	"context"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/binchencoder/letsgo/testing/mocks/etcd"
	pb "github.com/binchencoder/skylb-api/proto"
)

const (
	keyTestService   = "/registry/services/endpoints/default/test-service"
	epKeyTestService = "/registry/services/endpoints/default/test-service/172.0.0.100_8080"
)

func TestCalculateKey(t *testing.T) {
	Convey("Calculate etcd key from namespace and service name", t, func() {
		eh := endpointsHub{}
		key := eh.calculateKey("default", "test-service")
		So(key, ShouldEqual, keyTestService)
	})
}

func TestCalculateEndpointKey(t *testing.T) {
	Convey("Calculate etcd endpoint key from namespace, service name, host and port", t, func() {
		eh := endpointsHub{}
		key := eh.calculateEndpointKey("default", "test-service", "172.0.0.100", 8080)
		So(key, ShouldEqual, epKeyTestService)
	})
}

func TestRefreshKey(t *testing.T) {
	Convey("Refresh the given etcd key", t, func() {
		ctx := context.Background()

		Convey("When everything is fine", func() {
			etcdcli := new(etcd.KeysAPIMock)
			etcdcli.On("Set", ctx, keyTestService, "", refreshOpts).Return(nil, nil)
			eh := endpointsHub{
				etcdCli: etcdcli,
			}

			err := eh.refreshKey(context.Background(), keyTestService)
			So(err, ShouldBeNil)
		})

		Convey("When etcd returns error", func() {
			etcdcli := new(etcd.KeysAPIMock)
			eh := endpointsHub{
				etcdCli: etcdcli,
			}

			mockErr := errors.New("mock-error")
			etcdcli.On("Set", ctx, keyTestService, "", refreshOpts).Return(nil, mockErr)

			err := eh.refreshKey(context.Background(), keyTestService)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, mockErr)
		})
	})
}

func TestSetKey(t *testing.T) {
	Convey("Set the given etcd key sucessfully", t, func() {
		ctx := context.Background()
		spec := pb.ServiceSpec{
			Namespace:   "default",
			ServiceName: "test-service",
			PortName:    "grpc",
		}

		Convey("When everything is fine", func() {
			etcdcli := new(etcd.KeysAPIMock)
			eh := endpointsHub{
				etcdCli: etcdcli,
			}

			expectedVal := "{\"metadata\":{\"name\":\"172.0.0.100:8080\",\"namespace\":\"default\",\"creationTimestamp\":null},\"subsets\":[{\"addresses\":[{\"ip\":\"172.0.0.100\",\"targetRef\":{\"kind\":\"Pod\",\"namespace\":\"default\"}}],\"ports\":[{\"name\":\"grpc\",\"port\":8080}]}]}"
			etcdcli.On("Set", ctx, keyTestService, expectedVal, setOpts).Return(nil, nil)

			err := eh.setKey(context.Background(), keyTestService, &spec, "172.0.0.100", 8080, 0)
			So(err, ShouldBeNil)
		})

		Convey("When etcd returns error", func() {
			etcdcli := new(etcd.KeysAPIMock)
			eh := endpointsHub{
				etcdCli: etcdcli,
			}

			mockErr := errors.New("mock-error")
			expectedVal := "{\"metadata\":{\"name\":\"172.0.0.100:8080\",\"namespace\":\"default\",\"creationTimestamp\":null},\"subsets\":[{\"addresses\":[{\"ip\":\"172.0.0.100\",\"targetRef\":{\"kind\":\"Pod\",\"namespace\":\"default\"}}],\"ports\":[{\"name\":\"grpc\",\"port\":8080}]}]}"
			etcdcli.On("Set", ctx, keyTestService, expectedVal, setOpts).Return(nil, mockErr)

			err := eh.setKey(context.Background(), keyTestService, &spec, "172.0.0.100", 8080, 0)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, mockErr)
		})
	})
}
