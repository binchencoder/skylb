package main

import (
	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"

	"github.com/binchencoder/skylb-api/util"
)

const (
	endpointsKeyPrefix = "/registry/services/endpoints"
)

var (
	delOpts   etcd.DeleteOptions
	getOpts   etcd.GetOptions
	setOpts   etcd.SetOptions
	watchOpts etcd.WatcherOptions
)

func init() {
	delOpts = etcd.DeleteOptions{
		Recursive: true,
	}
	getOpts = etcd.GetOptions{
		Recursive: true,
	}
	setOpts = etcd.SetOptions{}
	watchOpts = etcd.WatcherOptions{
		Recursive: true,
	}
}

func startBridge(fromEndpoints, toEndpoints []string) {
	for {
		glog.Infof("Connecting to FROM etcd cluster ...")
		fromCli := createEtcdClient(fromEndpoints)
		glog.Infof("Connecting to TO etcd cluster ...")
		toCli := createEtcdClient(toEndpoints)

		// Initial load.
		if err := initialCopy(fromCli, toCli); err != nil {
			glog.Errorf("Failed to load initial etcd values, %+v.", err)
			continue
		}

		// Watch for changes.
		if err := watchChanges(fromCli, toCli); err != nil {
			glog.Errorf("Failed to watch etcd changes, %+v.", err)
			continue
		}

		break
	}
}

func createEtcdClient(endpoints []string) etcd.KeysAPI {
	cli, err := etcd.New(etcd.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		glog.Fatalf("Failed to create etcd client, %v.", err)
	}
	if err = cli.Sync(context.Background()); err != nil {
		glog.Fatalf("Failed to sync etcd cluster: %v", err)
	}

	machines := cli.Endpoints()
	if len(machines) == 0 || len(machines[0]) == 0 {
		glog.Fatalln("No etcd machines found")
	}
	glog.Infoln("Found etcd machines:", machines)

	return etcd.NewKeysAPI(cli)
}

func initialCopy(fromCli, toCli etcd.KeysAPI) error {
	values, err := fromCli.Get(context.Background(), endpointsKeyPrefix, &getOpts)
	if err != nil {
		return err
	}

	walkNodes(values.Node, func(node *etcd.Node) {
		setKey(toCli, node.Key, node.Value)
	})

	return nil
}

func watchChanges(fromCli, toCli etcd.KeysAPI) error {
	w := fromCli.Watcher(endpointsKeyPrefix, &watchOpts)

	for {
		values, err := w.Next(context.Background())
		if err != nil {
			glog.Errorf("Failed to get next watch event, %v", err)
			return err
		}

		glog.Infof("Monitored etcd change %s", values.Action)

		switch values.Action {
		case util.ActionCompareAndDelete, util.ActionCompareAndSwap, util.ActionCreate, util.ActionSet:
			walkNodes(values.Node, func(node *etcd.Node) {
				setKey(toCli, node.Key, node.Value)
			})
		case util.ActionDelete, util.ActionExpire:
			walkNodes(values.PrevNode, func(node *etcd.Node) {
				deleteKey(toCli, node.Key)
			})
		}
	}
}

func walkNodes(node *etcd.Node, callback func(node *etcd.Node)) {
	if len(node.Nodes) > 0 {
		for _, n := range node.Nodes {
			walkNodes(n, callback)
		}
	}

	if !node.Dir {
		callback(node)
	}
}

func setKey(toCli etcd.KeysAPI, key string, value string) {
	_, err := toCli.Set(context.Background(), key, value, &setOpts)
	if err != nil {
		glog.Errorf("Failed to set key %s, TO etcd might be out sync with FROM etcd. %v", key, err)
	}
}

func deleteKey(toCli etcd.KeysAPI, key string) {
	_, err := toCli.Delete(context.Background(), key, &delOpts)
	if err != nil {
		glog.Errorf("Failed to set key %s, TO etcd might be out sync with FROM etcd. %v", key, err)
	}
}
