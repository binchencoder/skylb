package hub

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	api "k8s.io/api/core/v1"

	pb "binchencoder.com/skylb-api/proto"
)

func TestSkypbEndpointsToMap(t *testing.T) {
	spec := pb.ServiceSpec{
		Namespace:   "default",
		ServiceName: "vexillary-demo",
		PortName:    "grpc-port",
	}
	eps := api.Endpoints{
		Subsets: []api.EndpointSubset{
			{
				Addresses: []api.EndpointAddress{
					{
						IP: "192.168.1.1",
					},
					{
						IP: "192.168.1.2",
					},
					{
						IP: "192.168.1.3",
					},
				},
				Ports: []api.EndpointPort{
					{
						Name: "grpc-port",
						Port: 8080,
					},
				},
			},
		},
	}

	Convey("When converting SkyLB endpoints proto to map", t, func() {
		epsMap := skypbEndpointsToMap(&spec, &eps)
		So(epsMap, ShouldHaveLength, 3)

		for _, ip := range []string{
			"192.168.1.1",
			"192.168.1.2",
			"192.168.1.3",
		} {
			Convey(fmt.Sprintf("%s should be included in the map", ip), func() {
				v, ok := epsMap[fmt.Sprintf("%s:8080", ip)]
				So(ok, ShouldBeTrue)
				So(v.IP, ShouldEqual, ip)
				So(v.Port, ShouldEqual, 8080)
			})
		}
	})
}
