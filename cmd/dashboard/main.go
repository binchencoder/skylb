package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kataras/iris"

	"github.com/binchencoder/letsgo"
	"github.com/binchencoder/skylb/dashboard"
	"github.com/binchencoder/skylb/dashboard/db"
)

const (
	keyFile  = "key.pem"
	certFile = "cert.pem"
)

var (
	hostPort  = flag.String("host-port", ":8050", "The web server host:port")
	certDir   = flag.String("cert-dir", "/opt/skylb-dashboard/certs", "The directory holding certification files")
	staticDir = flag.String("static-dir", "/opt/skylb-dashboard/static", "The directory holding static files")
	debugMode = flag.Bool("debug-mode", false, "Enable the debug mode which allows dummy login")
	etcdEps   = flag.String("etcd-endpoints", "http://localhost:2379", "The comma separated ETCD endpoints, e.g., http://etcd1:2379,http://etcd2:2379")
)

func usage() {
	fmt.Println(`Jingoal SkyLB dashboard.

Usage:
	skylb-dashboard [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func checkFlags() {
	if *certDir == "" {
		fmt.Println("Flag cert-dir has to be provided.")
		os.Exit(2)
	}
	if *staticDir == "" {
		fmt.Println("Flag static-dir has to be provided.")
		os.Exit(2)
	}
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))
	checkFlags()

	iris.Config.Gzip = true

	iris.StaticServe(*staticDir, "/static")

	db.Init(*debugMode)
	dashboard.Init(*debugMode, *staticDir, *etcdEps)

	fmt.Printf("Start SkyLB dashboard at https://%s.\n", *hostPort)
	iris.ListenTLS(*hostPort, getFilePath(*certDir, certFile), getFilePath(*certDir, keyFile))
}

func getFilePath(dir, file string) string {
	path := filepath.Join(dir, file)
	if _, err := os.Stat(path); err != nil {
		if err == os.ErrNotExist {
			fmt.Println("Can not find cert file,", err)
			os.Exit(2)
		}
		fmt.Println("Error,", err)
		os.Exit(2)
	}
	return path
}
