package hub

import (
	"context"
	"fmt"
	"testing"

	etcdcli "github.com/coreos/etcd/client"
	. "github.com/smartystreets/goconvey/convey"

	"binchencoder.com/letsgo/testing/mocks/etcd"
	pb "binchencoder.com/skylb-api/proto"
)

const (
	keyService1 = "/registry/services/endpoints/default/service1"
)

func TestAddObserver(t *testing.T) {
	specs := []*pb.ServiceSpec{
		{
			Namespace:   namespace,
			ServiceName: serviceName,
			PortName:    portName,
		},
	}

	Convey("Add a service discovery observer", t, func() {
		ctx := context.Background()

		Convey("When everything is fine", func() {
			resp := etcdcli.Response{
				Node: &etcdcli.Node{
					Key:           keyService1,
					CreatedIndex:  21,
					ModifiedIndex: 21,
					TTL:           0,
					Nodes: []*etcdcli.Node{
						{
							Key:           keyService1 + "/172.0.10.1_8080",
							CreatedIndex:  24072,
							ModifiedIndex: 24072,
							TTL:           8,
							Value:         `{"metadata":{"name":"172.0.10.1:8080","namespace":"default"},"subsets":[{"addresses":[{"ip":"172.0.10.1","targetRef":{"kind":"Pod","namespace":"default"}}],"ports":[{"name":"port","port":8080}]}]}`,
						},
						{
							Key:           keyService1 + "/172.0.10.2_8080",
							CreatedIndex:  24073,
							ModifiedIndex: 24073,
							TTL:           8,
							Value:         `{"metadata":{"name":"172.0.10.2:8080","namespace":"default"},"subsets":[{"addresses":[{"ip":"172.0.10.2","targetRef":{"kind":"Pod","namespace":"default"}}],"ports":[{"name":"port","port":8080}]}]}`,
						},
						{
							Key:           keyService1 + "/172.0.10.3_8080",
							CreatedIndex:  24076,
							ModifiedIndex: 24076,
							TTL:           10,
							Value:         `{"metadata":{"name":"172.0.10.3:8080","namespace":"default"},"subsets":[{"addresses":[{"ip":"172.0.10.3","targetRef":{"kind":"Pod","namespace":"default"}}],"ports":[{"name":"port","port":8080}]}]}`,
						},
					},
				},
			}
			etcdMock := new(etcd.KeysAPIMock)
			etcdMock.On("Get", ctx, keyService1, &getOpts).Return(&resp, nil)
			eh := endpointsHub{
				etcdCli:  etcdMock,
				services: serviceMap{},
			}

			ch, err := eh.AddObserver(specs, "192.168.0.1:8000", true /* resolveFull */)
			So(ch, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(eh.services, ShouldContainKey, keyService1)

			so := eh.services[keyService1]
			So(so, ShouldNotBeNil)
			So(len(so.observers), ShouldEqual, 1)
			So(len(so.endpoints), ShouldEqual, 3)
			So(so.endpoints, ShouldContainKey, "172.0.10.1:8080")
			So(so.endpoints["172.0.10.1:8080"].IP, ShouldEqual, "172.0.10.1")
			So(so.endpoints["172.0.10.1:8080"].Port, ShouldEqual, 8080)
			So(so.endpoints, ShouldContainKey, "172.0.10.2:8080")
			So(so.endpoints["172.0.10.2:8080"].IP, ShouldEqual, "172.0.10.2")
			So(so.endpoints["172.0.10.2:8080"].Port, ShouldEqual, 8080)
			So(so.endpoints, ShouldContainKey, "172.0.10.3:8080")
			So(so.endpoints["172.0.10.3:8080"].IP, ShouldEqual, "172.0.10.3")
			So(so.endpoints["172.0.10.3:8080"].Port, ShouldEqual, 8080)
			fmt.Println(so)

			co := so.observers[0]
			So(co, ShouldNotBeNil)
			So(co.clientAddr, ShouldEqual, "192.168.0.1:8000")
		})
	})
}

func TestRemoveObserver(t *testing.T) {
	spec := pb.ServiceSpec{
		Namespace:   namespace,
		ServiceName: serviceName,
		PortName:    portName,
	}

	observers := []*clientObject{}
	for _, cip := range []string{
		"192.168.0.12:33486",
		"192.168.0.10:30121",
	} {
		observers = append(observers, &clientObject{
			spec:       &spec,
			clientAddr: cip,
			stopCh:     make(chan struct{}),
		})
	}

	so := serviceObject{
		spec:      &spec,
		endpoints: serviceEndpoints{},
		observers: observers,
	}

	eh := endpointsHub{
		services: serviceMap{},
	}
	key := eh.calculateKey(spec.Namespace, spec.ServiceName)
	eh.services[key] = &so

	Convey("Remove one observer from the endpoints hub (one still remaining)", t, func() {
		eh.RemoveObserver([]*pb.ServiceSpec{&spec}, "192.168.0.12:33486")
		So(eh.services, ShouldContainKey, key)
		sobj := eh.services[key]
		So(len(sobj.observers), ShouldEqual, 1)
	})
}

func TestRemoveObserverFromSlice(t *testing.T) {
	spec := pb.ServiceSpec{
		Namespace:   namespace,
		ServiceName: serviceName,
		PortName:    portName,
	}

	clientAddrs := []string{
		"192.168.0.12:33486",
		"192.168.0.10:30121",
		"192.168.0.18:35469",
		"192.168.1.33:30983",
		"192.168.1.54:31128",
	}

	observers := []*clientObject{}
	for _, cip := range clientAddrs {
		observers = append(observers, &clientObject{
			spec:       &spec,
			clientAddr: cip,
			stopCh:     make(chan struct{}),
		})
	}

	Convey("Remove observer from managed list", t, func() {
		Convey("when the observer does not exist in list", func() {
			obs := removeObserverFromSlice(observers, &spec, "192.168.0.12:10000")
			So(len(obs), ShouldEqual, 5)
		})
		Convey("when the observer exists in list", func() {
			obs := removeObserverFromSlice(observers, &spec, "192.168.1.33:30983")
			So(len(obs), ShouldEqual, 4)
		})
	})
}

func TestRemoveObserverFromSlice_withDupSpecs(t *testing.T) {
	spec := pb.ServiceSpec{
		Namespace:   namespace,
		ServiceName: serviceName,
		PortName:    portName,
	}

	clientAddrs := []string{
		"192.168.0.12:33486",
		"192.168.0.10:30121", // Be duplicated at the 5th element.
		"192.168.0.18:35469",
		"192.168.1.33:30983",
		"192.168.0.10:30121", // Keep the dup at end of the slice.
	}

	observers := []*clientObject{}
	for _, cip := range clientAddrs {
		observers = append(observers, &clientObject{
			spec:       &spec,
			clientAddr: cip,
			stopCh:     make(chan struct{}),
		})
	}

	Convey("Remove observer from managed list", t, func() {
		Convey("when removing the dup observer", func() {
			obs := removeObserverFromSlice(observers, &spec, "192.168.0.10:30121")
			So(len(obs), ShouldEqual, 3)
			for _, ip := range observers {
				So(ip, ShouldNotEqual, "192.168.0.10:30121")
			}
		})
	})
}
