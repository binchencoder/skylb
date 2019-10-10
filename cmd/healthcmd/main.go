package main

// Connects directly to service's grpc port and invoke grpc health api.

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	hpb "google.golang.org/grpc/health/grpc_health_v1"

	"binchencoder.com/letsgo"
	jt "binchencoder.com/letsgo/time"
)

var (
	serverAddr = flag.String("server-addr", "127.0.0.1:10000", "The grpc server address, with format host:port")
	timeout    = flag.Duration("timeout", 3*time.Second, "The grpc dial and call timeout")

	start time.Time
)

type PingResponse struct {
	hpb.HealthCheckResponse
	TimeInMs float64
}

func usage() {
	fmt.Println(`grpchealth: the grpc service health checker.

Usage:
	grpchealth [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithTimeout(*timeout))
	if err != nil {
		reportErr(err)
		os.Exit(2)
	}
	defer conn.Close()

	healthCli := hpb.NewHealthClient(conn)
	req := hpb.HealthCheckRequest{}
	start = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	hresp, err := healthCli.Check(ctx, &req, grpc.FailFast(false))
	// If service cannot be reached, grpc will print error to stdout,
	// see resetTransport() in google.golang.org/grpc/clientconn.go
	if nil != err {
		cancel()
		reportErr(err)
		os.Exit(2)
	}
	resp := &PingResponse{
		HealthCheckResponse: *hresp,
		TimeInMs:            jt.MillisecondSince(start),
	}
	b, err := json.Marshal(&resp)
	if nil != err {
		reportErr(err)
		os.Exit(2)
	}
	fmt.Printf("%s", b)
}

func reportErr(err error) {
	glog.Warningf("%#v", err)
	fmt.Print(genErrResp())
}

func genErrResp() string {
	hErrResp := &hpb.HealthCheckResponse{
		//FIXME pb omitempty causing "{}", or use own json
		Status: hpb.HealthCheckResponse_UNKNOWN,
	}
	resp := &PingResponse{
		HealthCheckResponse: *hErrResp,
		TimeInMs:            jt.MillisecondSince(start),
	}
	b, err := json.Marshal(&resp)
	if nil != err {
		panic(err)
	}
	return string(b)
}
