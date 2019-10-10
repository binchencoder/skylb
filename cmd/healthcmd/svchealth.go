package main

// Calls service's grpc health api via skylb.
//
// To check service availability:
// 	Specify "--target-service=xxx-service-name" besides
// 	correct --skylb-endpoints value.
//
// To check skylb server functionality:
//	Specify "--skylb-endpoints=<skylb Endpoint>" respectively
// 	for each skylb server.

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"binchencoder.com/letsgo"
	"binchencoder.com/letsgo/service/naming"
	skylbserver "binchencoder.com/skylb-api/server"
	"binchencoder.com/skylb/cmd/grpchealth"
	"binchencoder.com/skylb/hub"
	vexpb "binchencoder.com/ease-gateway/proto/data"
)

var (
	defaultSvcName, _ = naming.ServiceIdToName(vexpb.ServiceId_DORY_SERVICE)
	targetServiceName = flag.String("target-service", defaultSvcName,
		"The target service name. "+
			"If --etcd-endpoints is specified, --target-service will be ignored, "+
			"and target service names will be extracted from etcd.")
	checkCount = flag.Int("check-count", 1, "Check how many times.")

	selfService = flag.Bool("self-service", false, "Register self as service and test with self. "+
		"If this flag is true, --target-service will be ignored.")
)

func usage() {
	fmt.Println(`svchealth.

Usage:
	svchealth [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))
	flag.Set("v", "3")
	flag.Set("logtostderr", "true")

	if *selfService {
		// Test self service.
		// Register self as service.
		port := 13131
		sid := vexpb.ServiceId_SKYLB_FRONTEND
		go func() {
			skylbserver.Register(sid, "grpc", port)
			skylbserver.Start(fmt.Sprintf(":%d", port), func(s *grpc.Server) error {
				// No need to hpb.RegisterHealthServer(s, health.NewServer()),
				// since this is already done in skylb api.
				return nil
			})
		}()
		time.Sleep(time.Second)
		sids, _ := naming.ServiceIdToName(sid)
		grpchealth.CheckService(sid, sids, *checkCount)
	} else {
		// Test another real service.
		etcdCli := hub.CreateEtcdClient(*hub.EtcdEndpoints, false)
		if nil == etcdCli {
			// Check a single service.
			serviceId, err := naming.ServiceNameToId(*targetServiceName)
			if nil != err {
				glog.Error(err)
				return
			}
			grpchealth.CheckService(serviceId, *targetServiceName, *checkCount)
		} else {
			// Check all services currently registered in etcd.
			grpchealth.CheckAllServices(etcdCli)
		}
	}
}
