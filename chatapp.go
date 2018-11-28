package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func exit_message(message string) bool {
	return strings.TrimSpace(strings.ToLower(message)) == "quit"
}

func server_setup(port string) net.Conn {
	l, _ := net.Listen("tcp", ":"+port)
	connection, _ := l.Accept()
	return connection
}

func client_setup(host string, port string) net.Conn {
	connection, _ := net.Dial("tcp", host+":"+port)
	return connection
}

func listen_loop(connection net.Conn, quit chan bool) {
	reader := bufio.NewReader(connection)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if len(message) > 0 {
			fmt.Println("Received message: " + message)
		}
		if exit_message(message) {
			quit <- true
		}
	}
}

func send_loop(connection net.Conn, quit chan bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		message, _ := reader.ReadString('\n')
		fmt.Fprintf(connection, message+"\n")
		if exit_message(message) {
			quit <- true
		}
	}
}

func main() {
	args := os.Args[1:]
	node_type := args[0]
	host := args[1]
	port := args[2]

	var connection net.Conn
	if node_type == "client" {
		connection = client_setup(host, port)
	} else if node_type == "server" {
		connection = server_setup(port)
	}

	fmt.Println("Future Gadget Laboratory Chat Room")
	fmt.Println("----------------------------------")

	quit := make(chan bool)

	go listen_loop(connection, quit)
	go send_loop(connection, quit)

	<-quit
	fmt.Println("Terminating")
}
