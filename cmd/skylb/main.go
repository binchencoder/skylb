package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	hpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/binchencoder/letsgo"
	jgrpc "github.com/binchencoder/letsgo/grpc"
	lmetrics "github.com/binchencoder/letsgo/metrics"
	"github.com/binchencoder/letsgo/runtime/pprof"
	"github.com/binchencoder/skylb-api/metrics"
	pb "github.com/binchencoder/skylb-api/proto"
	"github.com/binchencoder/skylb/rpc"
)

var (
	hostPort   = flag.String("host-port", ":1900", "The gRPC server host:port")
	scrapeAddr = flag.String("scrape-addr", ":1920", "The address to listen on for HTTP requests.")
)

func usage() {
	fmt.Println(`SkyLB: the external gRPC Load balancer.

Usage:
	skylb [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	lis, err := net.Listen("tcp", *hostPort)
	if err != nil {
		glog.Fatalf("failed to listen: %v\n", err)
	}

	m := cmux.New(lis)

	// Match connections in order: first grpc, then HTTP.
	grpcl := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpl := m.Match(cmux.HTTP1Fast())

	go startHTTPServer(httpl)
	go startGrpcServer(grpcl)

	m.Serve()
}

func startGrpcServer(grpcl net.Listener) {
	glog.Infof("Starting SkyLB ...")

	streamMetricsInt := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		spec := &pb.ServiceSpec{ServiceName: "skylb"}
		return metrics.StreamServerInterceptor(spec, srv, ss, info, handler)
	}
	streamIncepts := make([]grpc.StreamServerInterceptor, 0, 2)
	streamIncepts = append(streamIncepts, jgrpc.StreamRecoverServerInterceptor, streamMetricsInt)

	unaryInt := grpc.UnaryInterceptor(jgrpc.UnaryRecoverServerInterceptor)

	s := grpc.NewServer(unaryInt, grpc.StreamInterceptor(jgrpc.ChainStreamServer(streamIncepts...)))

	pb.RegisterSkylbServer(s, rpc.NewSkylbServer())
	hpb.RegisterHealthServer(s, health.NewServer())

	glog.Infof("SkyLB grpc service started on %s.\n", *hostPort)

	if err := s.Serve(grpcl); err != nil {
		panic(err)
	}
}

func startHTTPServer(httpl net.Listener) {
	lmetrics.EnablePrometheus(http.DefaultServeMux)
	pprof.EnablePprof(http.DefaultServeMux)
	if err := http.Serve(httpl, nil); err != nil {
		glog.Fatalf("Failed to start prometheus server: %v", err)
	}
}
