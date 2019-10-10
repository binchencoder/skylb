package dashboard

import (
	"context"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"

	"binchencoder.com/letsgo/strings"
	"binchencoder.com/skylb-api/prefix"
)

var (
	etcdCli etcd.KeysAPI
)

func initEtcdClient(etcdEps string) {
	etcdCli = createEtcdClient(etcdEps)
	prefix.Init(etcdCli)
}

func createEtcdClient(etcdEps string) etcd.KeysAPI {
	eps := strings.CsvToSlice(etcdEps)
	if len(eps) == 0 {
		glog.Fatalln("Flag --etcd-endpoints is required.")
	}

	glog.Infof("Use etcd endpoints %s.", eps)

	var cli etcd.Client
	for {
		var err error
		if cli, err = etcd.New(etcd.Config{
			Endpoints: eps,
		}); err != nil {
			glog.Errorf("Failed to create etcd client, %v. Will retry after one second.", err)
			time.Sleep(time.Second)
			continue
		}
		if err = cli.Sync(context.Background()); err != nil {
			glog.Errorf("Failed to sync cluster: %v. Will retry after one second.", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	machines := cli.Endpoints()
	if len(machines) == 0 || len(machines[0]) == 0 {
		glog.Fatalln("No etcd machines found")
	}
	glog.Infoln("Found etcd machines:", machines)

	return etcd.NewKeysAPI(cli)
}
