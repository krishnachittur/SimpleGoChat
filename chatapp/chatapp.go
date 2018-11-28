package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"./client"
	"./server"
)

// Endpoint is a client or server endpoint
type Endpoint interface {
	Setup(host string, port string)
	RunLoop(quit chan os.Signal)
}

type config struct {
	NodeType string
	HostIP   string
	HostPort string
}

func (c config) isServer() bool {
	return strings.ToLower(c.NodeType) == "server"
}

func (c config) isClient() bool {
	return strings.ToLower(c.NodeType) == "client"
}

func getConfig() (p config) {
	defer flag.Parse()
	flag.StringVar(&p.NodeType, "node_type", "client", "'client' or 'server'")
	flag.StringVar(&p.HostIP, "host_ip", "localhost", "IP of host application.")
	flag.StringVar(&p.HostPort, "host_port", "5050", "Port of host application.")
	return
}

func sigTermChannel() chan os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	return quit
}

func main() {
	cfg := getConfig()

	var node Endpoint
	if cfg.isClient() {
		node = &Client.Client{}
	} else if cfg.isServer() {
		node = &Server.Server{}
	} else {
		fmt.Println("Error: `node_type` must be client or server")
		return
	}

	node.Setup(cfg.HostIP, cfg.HostPort)
	node.RunLoop(sigTermChannel())
}
