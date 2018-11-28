package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Server struct {
	connection net.Conn
	//TODO
}

func (this *Server) setup(host string, port string) {
	l, _ := net.Listen("tcp", ":"+port)
	this.connection, _ = l.Accept()
}

func (this *Server) listen_loop(quit chan bool) {
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

func (this *Server) send_loop(quit chan bool) {
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

func (this *Server) run_loop(quit chan bool) {
	go this.listen_loop(quit)
	go this.send_loop(quit)
}
