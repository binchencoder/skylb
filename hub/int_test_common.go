package hub

// fuyc: has to make this file built into go_default_library so as to be
// accessible by tests under other dirs (e.g. cmd/webserver/svclist/)

import (
	"flag"

	etcd "github.com/coreos/etcd/client"
)

// GetTestEtcdClient returns etcd client initialized either by flag
// -etcd-endpoints or by hard-coded etcd endpoint.
func GetTestEtcdClient() etcd.KeysAPI {
	// Allow taking arg from CLI with "-etcd-endpoints=http://host:port"
	flag.Parse()
	if "" == *EtcdEndpoints {
		ep := "http://192.168.100.190:2379"
		EtcdEndpoints = &ep
	}
	return CreateEtcdClient(*EtcdEndpoints, true)
}
