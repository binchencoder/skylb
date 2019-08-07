package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	keyTestService   = "/registry/services/endpoints/default/test-service"
	epKeyTestService = "/registry/services/endpoints/default/test-service/172.0.0.100_8080"
)

func TestCalculateKey(t *testing.T) {
	Convey("Calculate etcd key from namespace and service name", t, func() {
		key := CalculateKey("default", "test-service")
		So(key, ShouldEqual, keyTestService)
	})
}

func TestCalculateEndpointKey(t *testing.T) {
	Convey("Calculate etcd endpoint key from namespace, service name, host and port", t, func() {
		key := CalculateEndpointKey("default", "test-service", "172.0.0.100", 8080)
		So(key, ShouldEqual, epKeyTestService)
	})
}
