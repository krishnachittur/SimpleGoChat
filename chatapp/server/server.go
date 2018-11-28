package Server

import (
	"bufio"
	"net"
	"os"
	"sync"
	"errors"
	"../csprotocol"
)

// Server implements the Endpoint interface. It services clients.
type Server struct {
	networkListener net.Listener

	roomnameToRoom   map[string]*Chatroom
	roomnameToRoomLock sync.Mutex
}

// Setup sets up the connection to listen on the supplied port
func (server *Server) Setup(host string, port string) {
	server.networkListener, _ = net.Listen("tcp", ":"+port)
}

// Determines whether or not to let requestingClient either a) Make a new chatroom or b) join an existing chatroom
func (server *Server) resolveChatroomReq(requestingClient *ClientConnection, chatroomReq csprotocol.ChatroomReq) error {
	server.roomnameToRoomLock.Lock()
	defer server.roomnameToRoomLock.Unlock()

	chatroom, ok := server.roomnameToRoom[chatroomReq.ChatroomID]

	if chatroomReq.IsNewChatroom {
		if ok {
			return errors.New("trying to create a chatroom that already exists")
		}
		newRoom := &Chatroom{ID: chatroomReq.ChatroomID, password: chatroomReq.ChatroomPassword}
		newRoom.clients = append(newRoom.clients, requestingClient)
		server.roomnameToRoom[chatroomReq.ChatroomID] = newRoom
	} else {
		if !ok {
			return errors.New("trying to join a chatroom doesn't exist")
		}
		chatroom.clientsLock.Lock()
		defer chatroom.clientsLock.Unlock()
		chatroom.clients = append(chatroom.clients, requestingClient)
	}
	return nil
}

// Broadcasts a message to everyone in the requestingClient's chatroom.
func (server *Server) resolveMessageBroadcastReq(requestingClient *ClientConnection, msgBcstRq csprotocol.MessageBroadcastReq) error {
	server.roomnameToRoomLock.Lock()
	chatroom, _ := server.roomnameToRoom[requestingClient.Roomname]
	server.roomnameToRoomLock.Unlock()

	msgNtf := csprotocol.MessageNotification{Message: msgBcstRq.Message, ClientID: requestingClient.ID}
	chatroom.broadcastToAllExcept(requestingClient, msgNtf)
	return nil
}

// Determines client intent then forever listens for messages (i.e. Message Broadcast Requests) from the client and forwards them.
func (server *Server) handleClient(clientConn net.Conn) {
	client := &ClientConnection{
		Connection: clientConn,
		ConnectionReader: bufio.NewReader(clientConn),
	}

	// Determine the client's identity.
	client.resolveIdentityReq()

	// Determine the chatroom the client wants to join.
	client.resolveRoomReq(
		func (rq *ClientConnection, crq csprotocol.ChatroomReq) error {
			return server.resolveChatroomReq(rq, crq)
		})

	for {
		// Wait for a Message Broadcast request from the client.
		// Then, broadcast the message in the chatroom.
		client.resolveMessageBroadcastReq(
			func (rq *ClientConnection, mrq csprotocol.MessageBroadcastReq) error {
				return server.resolveMessageBroadcastReq(rq, mrq)
			})
	}
}

// Forever listen for new connections.
func (server *Server) registerNewConnections() {
	for {
		newConnection, _ := server.networkListener.Accept()
		go server.handleClient(newConnection)
	}
}

func (server *Server) terminate() {
	// TODO
}

// RunLoop spawns a go thread that listens for new client connections.
func (server *Server) RunLoop(quit chan os.Signal) {
	go server.registerNewConnections()
	<-quit
	server.terminate()
}
