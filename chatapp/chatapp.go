package main

import (
	"fmt"
	"os"
	"strings"
)

//client or server endpoint
type Endpoint interface {
	setup(host string, port string)
	run_loop(quit chan bool)
}

func exit_message(message string) bool {
	return strings.TrimSpace(strings.ToLower(message)) == "\\quit"
}

func print_help() {
	fmt.Println("Usage: <program> [\"client\"|\"server\"] <host> <port>")
	fmt.Println("The <host> field is ignored if this program is run as a server.")
}

func main() {
	if len(os.Args) != 4 {
		print_help()
		return
	}
	args := os.Args[1:]
	node_type := args[0]
	host := args[1]
	port := args[2]

	var node Endpoint
	if node_type == "client" {
		node = &Client{}
	} else if node_type == "server" {
		node = &Server{}
	} else {
		print_help()
		return
	}

	node.setup(host, port)

	fmt.Println("Future Gadget Laboratory Chat Room")
	fmt.Println("----------------------------------")

	quit := make(chan bool)

	go node.run_loop(quit)

	<-quit
	fmt.Println("Terminating")
}
