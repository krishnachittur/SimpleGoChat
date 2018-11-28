package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	connection net.Conn
}

func (this *Client) setup(host string, port string) {
	this.connection, _ = net.Dial("tcp", host+":"+port)
}

func (this *Client) listen_loop(quit chan bool) {
	reader := bufio.NewReader(this.connection)
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

func (this *Client) send_loop(quit chan bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		message, _ := reader.ReadString('\n')
		fmt.Fprintf(this.connection, message+"\n")
		if exit_message(message) {
			quit <- true
		}
	}
}

func (this *Client) run_loop(quit chan bool) {
	go this.listen_loop(quit)
	go this.send_loop(quit)
}