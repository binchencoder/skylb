package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	pb "github.com/binchencoder/skylb-api/proto"
)

func TestServiceObject(t *testing.T) {
	Convey("Create a new service object and add an observer", t, func() {
		spec := pb.ServiceSpec{
			Namespace:   "default",
			ServiceName: "test-service",
			PortName:    "grpc",
		}

		so := NewServiceObject(&spec, nil)
		realSo := so.(*serviceObject)

		Convey("Get spec", func() {
			So(so.Spec(), ShouldResemble, &spec)
		})

		Convey("Set endpoints", func() {
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
			so.SetEndpoints(&se)
			So(realSo.endpoints, ShouldResemble, &se)
			So(realSo.observers, ShouldHaveLength, 0)

			Convey("Add an observer, notification should be sent", func() {
				co := ClientObserverMock{}
				co.On("ClientAddr").Return("192.168.0.1:33254")
				co.On("Notify", &se)
				so.AddObserver(&co)
				So(realSo.observers, ShouldHaveLength, 1)

				Convey("After adding the observer, the service object should hold one observer", func() {
					So(realSo.observers, ShouldHaveLength, 1)
				})
			})
		})
	})

	Convey("Create a new service object, add an observer, then remove it", t, func() {
		spec := pb.ServiceSpec{
			Namespace:   "default",
			ServiceName: "test-service",
			PortName:    "grpc",
		}

		so := NewServiceObject(&spec, nil)
		realSo := so.(*serviceObject)

		Convey("Set endpoints", func() {
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
			so.SetEndpoints(&se)
			So(realSo.endpoints, ShouldResemble, &se)
			So(realSo.observers, ShouldHaveLength, 0)

			Convey("Add an observer, then remove it", func() {
				co := ClientObserverMock{}
				co.On("ClientAddr").Return("192.168.0.1:33254")
				co.On("Close")
				co.On("Notify", &se)
				co.On("Spec").Return(&spec)
				so.AddObserver(&co)
				So(realSo.observers, ShouldHaveLength, 1)
				so.RemoveObservers("192.168.0.1:33254")

				Convey("After adding then removing the observer, the service object should hold zero observer", func() {
					So(realSo.observers, ShouldHaveLength, 0)
				})
			})
		})
	})

	// The old SkyLB Java API might send multiple requests with the same gRPC
	// connection thus for SkyLB server they all have the same client addr.
	//
	// This test case makes sure that RemoveObservers() will correctly remove
	// all observers with the same client addr.
	Convey("Remove dup observers from a service object", t, func() {
		spec := pb.ServiceSpec{
			Namespace:   "default",
			ServiceName: "test-service",
			PortName:    "grpc",
		}

		so := NewServiceObject(&spec, nil)
		realSo := so.(*serviceObject)
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
		so.SetEndpoints(&se)

		for i := 0; i < 5; i++ {
			co := ClientObserverMock{}
			co.On("ClientAddr").Return("192.168.0.1:33254")
			co.On("Close")
			co.On("Notify", &se)
			co.On("Spec").Return(&spec)
			so.AddObserver(&co)
		}
		So(realSo.observers, ShouldHaveLength, 5)

		Convey("Remove observers with the same ids, the service object should hold zero observer", func() {
			so.RemoveObservers("192.168.0.1:33254")
			So(realSo.observers, ShouldHaveLength, 0)
		})
	})
}
