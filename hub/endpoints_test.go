package hub

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	api "k8s.io/api/core/v1"

	pb "binchencoder.com/skylb-api/proto"
)

const (
	namespace   = "default"
	serviceName = "service1"
	portName    = "port"
	port        = 8080
)

func TestDiffEndpoints(t *testing.T) {
	spec := &pb.ServiceSpec{
		Namespace:   namespace,
		ServiceName: serviceName,
		PortName:    portName,
	}

	last := serviceEndpoints{
		"192.168.1.1:8080": ServiceEndpoint{
			IP:   "192.168.1.1",
			Port: 8080,
		},
		"192.168.1.2:8080": ServiceEndpoint{
			IP:   "192.168.1.2",
			Port: 8080,
		},
		"192.168.1.3:8080": ServiceEndpoint{
			IP:   "192.168.1.3",
			Port: 8080,
		},
	}

	now := serviceEndpoints{
		"192.168.1.2:8080": ServiceEndpoint{
			IP:   "192.168.1.2",
			Port: 8080,
		},
		"192.168.1.3:8080": ServiceEndpoint{
			IP:   "192.168.1.3",
			Port: 8080,
		},
		"192.168.1.4:8080": ServiceEndpoint{
			IP:   "192.168.1.4",
			Port: 8080,
		},
	}

	Convey("Calculate endpoints diff", t, func() {
		diff := diffEndpoints(spec, last, now)
		So(diff, ShouldNotBeNil)
		So(diff.Spec, ShouldNotBeNil)
		So(diff.Spec.Namespace, ShouldEqual, namespace)
		So(diff.Spec.ServiceName, ShouldEqual, serviceName)
		So(diff.Spec.PortName, ShouldEqual, portName)
		So(len(diff.InstEndpoints), ShouldEqual, 2)

		Convey("verify the 1st endpoint", func() {
			ep := diff.InstEndpoints[0]
			So(ep.Op, ShouldEqual, pb.Operation_Delete)
			So(ep.Host, ShouldEqual, "192.168.1.1")
			So(ep.Port, ShouldEqual, port)
		})

		Convey("verify the 2nd endpoint", func() {
			ep := diff.InstEndpoints[1]
			So(ep.Op, ShouldEqual, pb.Operation_Add)
			So(ep.Host, ShouldEqual, "192.168.1.4")
			So(ep.Port, ShouldEqual, port)
		})
	})
}

func TestFindPort(t *testing.T) {
	ports := []api.EndpointPort{
		{
			Name: "udp-port",
			Port: 5000,
		},
		{
			Name: "grpc-port",
			Port: 5050,
		},
		{
			Name: "web-port",
			Port: 5051,
		},
	}

	Convey("Find endpoint port", t, func() {
		Convey("find udp port", func() {
			port := findPort(ports, "udp-port")
			So(port, ShouldEqual, 5000)
		})
		Convey("find grpc port", func() {
			port := findPort(ports, "grpc-port")
			So(port, ShouldEqual, 5050)
		})
		Convey("find web port", func() {
			port := findPort(ports, "web-port")
			So(port, ShouldEqual, 5051)
		})
		Convey("find non existing port", func() {
			port := findPort(ports, "grpcport")
			So(port, ShouldEqual, 0)
		})
	})
}
