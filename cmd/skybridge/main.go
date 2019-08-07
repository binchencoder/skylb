// skybridge copies changes of service endpoints from one etcd cluster
// to another.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"

	"github.com/binchencoder/letsgo"
	"github.com/binchencoder/letsgo/strings"
)

var (
	fromEtcdEndpoints = flag.String("from-etcd-endpoints", "", "The FROM etcd endpoints")
	toEtcdEndpoints   = flag.String("to-etcd-endpoints", "", "The TO etcd endpoints")
)

func usage() {
	fmt.Println(`SkyBridge: the etcd bridge for SkyLB in Kubernetes.

Usage:
	skybridge [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func checkFlags() {
	if *fromEtcdEndpoints == "" {
		glog.Fatalf("Flag --from-etcd-endpoints is required.")
	}
	if *toEtcdEndpoints == "" {
		glog.Fatalf("Flag --to-etcd-endpoints is required.")
	}
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))
	checkFlags()

	glog.Infof("From etcd endpoints: %s", *fromEtcdEndpoints)
	glog.Infof("To etcd endpoints: %s", *toEtcdEndpoints)

	fromEndpoints := strings.CsvToSlice(*fromEtcdEndpoints)
	toEndpoints := strings.CsvToSlice(*toEtcdEndpoints)
	startBridge(fromEndpoints, toEndpoints)
}
