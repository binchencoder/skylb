package hub

import (
	api "k8s.io/api/core/v1"

	pb "binchencoder.com/skylb-api/proto"
)

type EndpointsUpdate struct {
	Id        int64
	Endpoints *pb.ServiceEndpoints
}

func diffEndpoints(spec *pb.ServiceSpec, last, now serviceEndpoints) *pb.ServiceEndpoints {
	// Found common entries.
	common := make(map[string]struct{})
	for k := range now {
		if _, ok := last[k]; ok {
			common[k] = struct{}{}
		}
	}

	eps := []*pb.InstanceEndpoint{}

	// Found endpoints to be removed from client.
	for k, v := range last {
		if _, ok := common[k]; !ok {
			ep := pb.InstanceEndpoint{
				Op:   pb.Operation_Delete,
				Host: v.IP,
				Port: v.Port,
			}
			eps = append(eps, &ep)
		}
	}

	// Found endpoints to be added for client.
	for k, v := range now {
		if _, ok := common[k]; !ok {
			ep := pb.InstanceEndpoint{
				Op:     pb.Operation_Add,
				Host:   v.IP,
				Port:   v.Port,
				Weight: v.Weight,
			}
			eps = append(eps, &ep)
		}
	}

	return &pb.ServiceEndpoints{
		Spec:          spec,
		InstEndpoints: eps,
	}
}

func findPort(ports []api.EndpointPort, portName string) int32 {
	for _, ep := range ports {
		if ep.Name == portName {
			return ep.Port
		}
	}
	return 0
}
