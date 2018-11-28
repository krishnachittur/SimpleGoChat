package Server

import (
	"bufio"
	"net"
	"os"
	"sync"
	"../csprotocol"
)

type Server struct {
	networkListener net.Listener

	roomnameToRoom   map[string]Chatroom
	roomnameToRoomLock sync.Mutex
}

func (server *Server) Setup(host string, port string) {
	server.networkListener, _ = net.Listen("tcp", ":"+port)
}

func (server *Server) resolveChatroomRequest(requestingClient ClientConnection, chatroomReq csprotocol.ChatroomReq) err error {
	server.roomnameToRoomLock.Lock()
	defer server.roomnameToRoomLock.Unlock()

	chatroom, ok := server.roomnameToRoom[chatroomReq.ChatroomID]

	if chatroomReq.IsNewChatroom {
		if ok {
			error "trying to create a chatroom that already exists"
			return
		}
		newRoom := Chatroom{ID: chatroomReq.ChatroomID, password: chatroomReq.ChatroomPassword}
		newRoom.clients = append(newRoom.clients, requestingClient)
		server.roomnameToRoom[chatroomReq.ChatroomID] = newRoom
		return true
	} else {
		if !ok {
			error "trying to join a chatroom doesn't exist"
			return false
		}
		chatroom.clientsLock.Lock()
		defer chatroom.clientsLock.Unlock()
		chatroom.clients = append(chatroom.clients, requestingClient)
		return true
	}
}

func (server *Server) resolveMessageBroadcastReq(requestingClient ClientConnection, msgBcstRq csprotocol.MessageBroadcastReq) bool {
	roomnameToRoomLock.Lock()
	chatroom, _ := server.roomnameToRoom[requestingClient.roomname]
	roomnameToRoomLock.Unlock()

	MessageNotification{Message: msgBcstRq, ClientID: requestingClient.ID}
	chatroom.broadcastToAllExcept(requestingClient, msgBcstRq)
}

func (server *Server) newClient(clientConn net.Conn) {
	client := ClientConnection{
		Connection: clientConn,
		ConnectionReader: bufio.NewReader(clientConn),
	}
	client.resolveIdentityReq()
	client.resolveRoomRequest(
		func (rq ClientConnection, crq csprotocol.ChatroomReq) error {
			return server.resolveChatroomRequest(rq, crq)
		})

	for {
		client.resolveMessageBroadcastReq(
			func (rq ClientConenction, mrq csprotocol.MessageBroadcastReq) error {
				return server.resolveMessageBroadcastRequest(rq, mrq)
			})
	}
}

func (server *Server) registerNewConnections() {
	for {
		newConnection, _ := server.networkListener.Accept()
		go server.newClient(newConnection)
	}
}

func (server *Server) terminate() {
	// TODO
}

func (server *Server) RunLoop(quit chan os.Signal) {
	go server.registerNewConnections()
	<-quit
	server.terminate()
}
