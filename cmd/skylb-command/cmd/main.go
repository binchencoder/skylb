package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/peterh/liner"
	api "k8s.io/api/core/v1"

	"github.com/binchencoder/letsgo"
)

const (
	prefix          = "/registry/services/endpoints"
	kind            = "Pod"
	defaultPortName = "grpc"
	defaultPromp    = ">> "
)

type service struct {
	name      string
	namespace string
	endpoints []string
}

var (
	etcdEndpoints = flag.String("etcd-endpoints", "", "Comma separated etcd endpoints")

	line     *liner.State
	commands = map[string]string{
		"?":        "Show this help",
		"exit":     "Exit the interactive shell",
		"help":     "Show this help",
		"quit":     "Quit the interactive shell",
		"add":      "Add a new instance for the current service",
		"ls":       "List all services or instances of a service, depending on the context",
		"new":      "Create a new service",
		"rm":       "Delete an instance for the current service",
		"portname": "Display or set port name",
		"select":   "Select a service to manage or reset to not manage any service",
	}
	cmds []string

	getOpts etcd.GetOptions
	setOpts etcd.SetOptions
	delOpts etcd.DeleteOptions

	portRegex = regexp.MustCompile("^\\d+$")

	services       = make([]*service, 0, 100)
	currentService *service

	portName = defaultPortName
	prompt   = defaultPromp
)

func init() {
	cmds = make([]string, 0, len(commands))
	for k := range commands {
		cmds = append(cmds, k)
	}
	sort.Sort(sort.StringSlice(cmds))

	getOpts = etcd.GetOptions{
		Recursive: true,
	}
	setOpts = etcd.SetOptions{
		PrevExist: etcd.PrevNoExist,
	}
	delOpts = etcd.DeleteOptions{}
}

func usage() {
	fmt.Println(
		`skylb-command.

Usage:
	skylb-command [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func checkFlags() {
	if *etcdEndpoints == "" {
		fmt.Println("flag --etcd-endpoints is required")
		os.Exit(2)
	}
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))
	checkFlags()

	line = createLiner()
	defer line.Close()
	defer saveLiner(line)

	cli, err := createEtcdClient()
	if err != nil {
		fmt.Printf("Error, %s.\n", err.Error())
		return
	}

	for {
		if cmd, err := line.Prompt(prompt); err == nil {
			cmd = strings.TrimSpace(cmd)
			if cmd == "" {
				continue
			}
			line.AppendHistory(cmd)

			switch cmd {
			case "?", "help":
				printHelp()
			case "exit", "quit":
				return
			case "ls":
				if currentService == nil {
					listServices(cli)
				} else {
					listInstances(cli)
				}
			case "add":
				fmt.Println("\tusage: add <host>:<port>")
			case "new":
				fmt.Println("\tusage: new name-space service-name")
			case "portname":
				fmt.Printf("\tCurrent port name: %s\n", portName)
			case "rm":
				fmt.Println("\tusage: rm <index>")
			case "select":
				fmt.Println("\t Reset. Not managing any service. To select a service to manage, run \"select <number>\"")
				prompt = defaultPromp
				currentService = nil
			default:
				if strings.HasPrefix(cmd, "add ") {
					addInstance(cli, cmd[4:])
				} else if strings.HasPrefix(cmd, "new ") {
					params := strings.Split(cmd[4:], " ")
					if len(params) != 2 {
						fmt.Println("\tTwo parameters are expected: name-space service-name")
						continue
					}
					addService(cli, params)
				} else if strings.HasPrefix(cmd, "rm ") {
					deleteInstance(cli, cmd[3:])
				} else if strings.HasPrefix(cmd, "select ") {
					selectService(cli, cmd[7:])
				} else if strings.HasPrefix(cmd, "portname ") {
					setPortname(cli, cmd[9:])
				} else {
					fmt.Println("\tUnknown command.")
				}
			}
		} else if err == liner.ErrPromptAborted || err == io.EOF {
			fmt.Println()
			break
		} else {
			fmt.Print("\tError reading line: ", err)
		}
	}
}

func printHelp() {
	fmt.Println("Commands:")
	for _, k := range cmds {
		fmt.Printf("  %12s: %s\n", k, commands[k])
	}
}

func addService(cli etcd.KeysAPI, params []string) {
	if currentService != nil {
		fmt.Printf("Current service is %s, please unselect it.\n", currentService)
		return
	}

	opts := etcd.SetOptions{
		Dir:       true,
		PrevExist: etcd.PrevNoExist,
	}

	key := path.Join(prefix, params[0], params[1])
	if _, err := cli.Set(context.Background(), key, "", &opts); err != nil {
		fmt.Printf("\tError, %s.\n", err.Error())
		return
	}
}

func addInstance(cli etcd.KeysAPI, param string) {
	if currentService == nil {
		fmt.Println("No service is selected.")
		return
	}

	param = strings.TrimSpace(param)
	parts := strings.Split(param, ":")
	if len(parts) != 2 {
		fmt.Printf("\tError, expect instance endpoint like <host>:<port>.\n")
	}

	host := strings.TrimSpace(parts[0])
	if len(host) == 0 {
		fmt.Println("\tError, valid host address is required.")
		return
	}

	port := strings.TrimSpace(parts[1])
	if !portRegex.MatchString(port) {
		fmt.Println("\tError, valid port is required.")
		return
	}

	key := path.Join(prefix, currentService.namespace, currentService.name, fmt.Sprintf("%s_%s", host, port))
	portNum, _ := strconv.Atoi(port)
	eps := api.Endpoints{
		Subsets: []api.EndpointSubset{
			{
				Addresses: []api.EndpointAddress{
					{
						IP: host,
						TargetRef: &api.ObjectReference{
							Kind:      kind,
							Namespace: currentService.namespace,
						},
					},
				},
				Ports: []api.EndpointPort{
					{
						Name: portName,
						Port: int32(portNum),
					},
				},
			},
		},
	}
	eps.Name = fmt.Sprintf("%s:%s", host, port)
	eps.Namespace = currentService.namespace
	b, err := json.Marshal(&eps)
	if err != nil {
		fmt.Printf("\tError, %s.\n", err.Error())
		return
	}

	if _, err := cli.Set(context.Background(), key, string(b), &setOpts); err != nil {
		fmt.Printf("\tError, %s.\n", err.Error())
		return
	}
	fmt.Println("\tDone.")
}

func deleteInstance(cli etcd.KeysAPI, param string) {
	if currentService == nil {
		fmt.Println("No service is selected.")
		return
	}

	param = strings.TrimSpace(param)
	idx, err := strconv.Atoi(param)
	if err != nil {
		fmt.Printf("\tError, %s.\n", err.Error())
		return
	}

	if idx < 0 || idx >= len(currentService.endpoints) {
		fmt.Println("\tIndex exceeds limit. Use \"ls\" to list all endpoints.")
		return
	}

	eps := currentService.endpoints[idx]
	key := path.Join(prefix, currentService.namespace, currentService.name, strings.Replace(eps, ":", "_", -1))

	// Make sure the endpoint is static (without TTL).
	resp, err := cli.Get(context.Background(), key, &getOpts)
	if err != nil {
		fmt.Printf("\tFailed to get instance %s doesn't exist: %v.\n", eps, err)
	}

	if resp.Node.TTL > 0 {
		fmt.Println("\tCan not remove an endpoint with TTL set.")
		return
	}

	prompt := fmt.Sprintf("\tAre you sure to delete endpoint %s? ([y]/n): ", currentService.endpoints[idx])
	input, err := line.Prompt(prompt)
	if err != nil {
		fmt.Printf("\tError, %s\n", err.Error())
		return
	}

	input = strings.TrimSpace(input)
	if !(input == "" || input == "y" || input == "Y") {
		return
	}

	if _, err := cli.Delete(context.Background(), key, &delOpts); err != nil {
		fmt.Printf("\tError, %s.\n", err.Error())
		return
	}
}

func listServices(cli etcd.KeysAPI) {
	resp, err := cli.Get(context.Background(), prefix, &getOpts)
	if err != nil {
		fmt.Printf("\tFailed to list services: %v.\n", err)
	}

	sort.Sort(nodeSlice(resp.Node.Nodes))

	services = make([]*service, 0, 100)
	cnt := 0
	for _, nsNode := range resp.Node.Nodes {
		namespace := filepath.Base(nsNode.Key)
		for _, serviceNode := range nsNode.Nodes {
			serviceName := filepath.Base(serviceNode.Key)
			svc := service{
				namespace: namespace,
				name:      serviceName,
				endpoints: make([]string, 0, 10),
			}

			services = append(services, &svc)

			fmt.Printf("\t%d: %s@%s\n", cnt, serviceName, namespace)
			cnt++
		}
	}
	if cnt > 0 {
		fmt.Println()
	}
	fmt.Printf("\tFound %d services in total.\n", cnt)
}

func listInstances(cli etcd.KeysAPI) {
	if currentService == nil {
		return
	}

	key := path.Join(prefix, currentService.namespace, currentService.name)
	resp, err := cli.Get(context.Background(), key, &getOpts)
	if err != nil {
		if e, ok := err.(etcd.Error); ok && e.Code == etcd.ErrorCodeKeyNotFound {
			fmt.Printf("\tService %s.%s absent, return empty list.\n", currentService.namespace, currentService.name)
		} else {
			fmt.Printf("\tError, %s.\n", err.Error())
		}
		return
	}

	sort.Sort(nodeSlice(resp.Node.Nodes))

	currentService.endpoints = make([]string, 0, 10)
	cnt := 0
	for _, node := range resp.Node.Nodes {
		eps := api.Endpoints{}
		if err := json.Unmarshal([]byte(node.Value), &eps); err != nil {
			fmt.Printf("\tError, %s.\n", err.Error())
			return
		}

		ttl := "(Static)"
		if node.TTL > 0 {
			ttl = fmt.Sprintf("(TTL: %ds)", node.TTL)
		}
		for _, sub := range eps.Subsets {
			for k, port := range sub.Ports {
				if port.Name == portName {
					fmt.Printf("\t%d: %s:%d %s\n", cnt, sub.Addresses[k].IP, port.Port, ttl)
					currentService.endpoints = append(currentService.endpoints, fmt.Sprintf("%s:%d", sub.Addresses[k].IP, port.Port))
					cnt++
					break
				}
			}
		}
	}
	if cnt > 0 {
		fmt.Println()
	}
	fmt.Printf("\tFound %d instances in total.\n", cnt)
}

func selectService(cli etcd.KeysAPI, svcNum string) {
	num, err := strconv.Atoi(svcNum)
	if err != nil {
		fmt.Printf("\tWrong format. Should be \"select <number>\", %v.\n", err)
		return
	}

	if num < 0 || num >= len(services) {
		fmt.Printf("\tIndex exceeds limit. Use \"ls\" to list all services.\n")
		return
	}

	currentService = services[num]
	prompt = fmt.Sprintf("%s %s", currentService.name, defaultPromp)
}

func setPortname(cli etcd.KeysAPI, pn string) {
	fmt.Printf("Set port name to %s.\n", pn)
	portName = pn
}

func createEtcdClient() (etcd.KeysAPI, error) {
	eps := strings.Split(*etcdEndpoints, ",")
	if *etcdEndpoints == "" || len(eps) == 0 {
		return nil, errors.New("flag --etcd-endpoints is required")
	}

	fmt.Printf("Use etcd endpoints %s.\n", eps)

	var cli etcd.Client
	for {
		var err error
		if cli, err = etcd.New(etcd.Config{
			Endpoints: eps,
		}); err != nil {
			fmt.Printf("Failed to create etcd client, %v. Will retry after one second.\n", err)
			time.Sleep(time.Second)
			continue
		}
		if err = cli.Sync(context.Background()); err != nil {
			fmt.Printf("Failed to sync cluster: %v. Will retry after one second.\n", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	machines := cli.Endpoints()
	if len(machines) == 0 || len(machines[0]) == 0 {
		fmt.Println("No etcd machines found")
		os.Exit(2)
	}
	fmt.Println("Found etcd machines:", machines)

	return etcd.NewKeysAPI(cli), nil
}
