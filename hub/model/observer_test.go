package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	pb "github.com/binchencoder/skylb-api/proto"
)

func TestClientObserver(t *testing.T) {
	Convey("Create a new service client observer", t, func() {
		spec := pb.ServiceSpec{
			Namespace:   "default",
			ServiceName: "test-service",
			PortName:    "grpc",
		}
		clientAddr := "192.168.0.1:33254"
		notifyCh := make(chan *pb.ServiceEndpoints, 1)

		co := NewClientObserver(&spec, clientAddr, notifyCh)
		realCo := co.(*clientObserver)

		Convey("Get clientAddr", func() {
			So(co.ClientAddr(), ShouldEqual, clientAddr)
		})

		Convey("Get spec", func() {
			So(co.Spec(), ShouldResemble, &spec)
		})

		Convey("Notify the observer of endpoint changes", func() {
			se := pb.ServiceEndpoints{
				Spec: &spec,
				InstEndpoints: []*pb.InstanceEndpoint{
					{
						Host: "172.10.0.100",
						Op:   pb.Operation_Add,
						Port: 8000,
					},
				},
			}
			co.Notify(&se)

			So(notifyCh, ShouldHaveLength, 1)
			So(<-notifyCh, ShouldResemble, &se)

			Convey("After closing the observer, notification should not be sent", func() {
				So(realCo.closed, ShouldBeFalse)
				co.Close()
				So(realCo.closed, ShouldBeTrue)
				co.Notify(&se)
				So(notifyCh, ShouldHaveLength, 0)
			})
		})
	})
}
