package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"

	"binchencoder.com/letsgo/strings"
)

var (
	etcdEps0 = flag.String("etcd-endpoints", "http://localhost:2379", "The comma separated ETCD endpoints, e.g., http://etcd1:2379,http://etcd2:2379")

	setOpts *etcd.SetOptions
)

func init() {
	setOpts = &etcd.SetOptions{
		PrevExist: etcd.PrevNoExist,
	}
}

func usage0() {
	fmt.Println(`Jingoal SkyLB ETCD loader. Only used for dev environment.

Usage:
	etcd-loader [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage0
	flag.Parse()

	cli := createEtcdClient(*etcdEps0)

	// For dashboard, we only care about the key.

	ctx := context.Background()
	for _, key := range serviceKeys {
		if _, err := cli.Set(ctx, key, "", setOpts); err != nil {
			fmt.Println(err)
		}
	}
	for _, key := range graphKeys {
		if _, err := cli.Set(ctx, key, "", setOpts); err != nil {
			fmt.Println(err)
		}
	}
	for _, key := range lameduckKeys {
		if _, err := cli.Set(ctx, key, "", setOpts); err != nil {
			fmt.Println(err)
		}
	}
}

func createEtcdClient(etcdEndpoints string) etcd.KeysAPI {
	eps := strings.CsvToSlice(etcdEndpoints)
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
