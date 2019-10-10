package hub

import (
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	api "k8s.io/api/core/v1"

	"binchencoder.com/skylb-api/prefix"
	pb "binchencoder.com/skylb-api/proto"
)

const (
	defaultKind = "Pod"
)

var (
	etcdKeyTtl = flag.Duration("etcd-key-ttl", 10*time.Second, "The etcd key TTL")

	refreshOpts *etcd.SetOptions
	setOpts     *etcd.SetOptions
)

func init() {
	refreshOpts = &etcd.SetOptions{
		TTL:     *etcdKeyTtl,
		Refresh: true,
	}
	setOpts = &etcd.SetOptions{
		TTL: *etcdKeyTtl,
	}
}

func calculateWeightKey(host string, port int32) string {
	return fmt.Sprintf("%s_%d_weight", host, port)
}

func (eh *endpointsHub) calculateKey(namespace, serviceName string) string {
	return path.Join(prefix.EndpointsKey, namespace, serviceName)
}

func (eh *endpointsHub) calculateEndpointKey(namespace, serviceName, host string, port int32) string {
	return path.Join(prefix.EndpointsKey, namespace, serviceName, fmt.Sprintf("%s_%d", host, port))
}

func (eh *endpointsHub) refreshKey(ctx context.Context, key string) error {
	glog.V(6).Infof("etcd set %#v -- %#v | %#v\n", key, "", refreshOpts)
	_, err := eh.etcdCli.Set(ctx, key, "", refreshOpts)
	return err
}

func (eh *endpointsHub) setKey(ctx context.Context, key string, spec *pb.ServiceSpec, host string, port, weight int32) error {
	eps := api.Endpoints{
		Subsets: []api.EndpointSubset{
			{
				Addresses: []api.EndpointAddress{
					{
						IP: host,
						TargetRef: &api.ObjectReference{
							Kind:      defaultKind,
							Namespace: spec.Namespace,
						},
					},
				},
				Ports: []api.EndpointPort{
					{
						Name: spec.PortName,
						Port: port,
					},
				},
			},
		},
	}
	if weight != 0 {
		eps.Labels = map[string]string{
			calculateWeightKey(host, port): fmt.Sprintf("%d", weight),
		}
	}
	eps.Name = fmt.Sprintf("%s:%d", host, port)
	eps.Namespace = spec.Namespace
	b, err := json.Marshal(&eps)
	if err != nil {
		return err
	}

	glog.V(6).Infof("etcd set %#v -- %#v | %#v\n", key, string(b), setOpts)
	_, err = eh.etcdCli.Set(ctx, key, string(b), setOpts)
	return err
}
