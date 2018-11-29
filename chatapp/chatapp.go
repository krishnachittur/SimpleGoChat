package main

import (
	"flag"
	"log"
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

func (c config) validateConfig() {
	if !c.isServer() && !c.isClient() {
		log.Fatal("The `node_type` flag is required! Please set it to either `client` or `server`.")
	} else if c.HostIP == "" {
		log.Fatal("The `host_ip` flag is required! If you're trying to run the client code, set it to the Host's IP." +
			"If you're trying to run the server code, set it to `localhost`")
	} else if c.HostPort == "" {
		log.Fatal("The `host_port` flag is required! If you're trying to run the client code, set it to the Host's Port." +
			"If you're trying to run the server code, set it to some unused port > 1024.")
	}
}

func getConfig() (p config) {
	defer flag.Parse()
	flag.StringVar(&p.NodeType, "node_type", "", "'client' or 'server'")
	flag.StringVar(&p.HostIP, "host_ip", "", "IP of host application.")
	flag.StringVar(&p.HostPort, "host_port", "", "Port of host application.")
	return
}

func sigTermChannel() chan os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	return quit
}

func main() {
	cfg := getConfig()
	cfg.validateConfig()

	var node Endpoint
	if cfg.isClient() {
		node = &Client.Client{}
	} else {
		node = &Server.Server{}
	}

	node.Setup(cfg.HostIP, cfg.HostPort)
	node.RunLoop(sigTermChannel())
}
