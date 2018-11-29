package Server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"../csprotocol"
)

// Server implements the Endpoint interface. It services clients.
type Server struct {
	networkListener net.Listener

	roomnameToRoom     map[string]*Chatroom
	roomnameToRoomLock *sync.Mutex
}

// Setup sets up the connection to listen on the supplied port
func (server *Server) Setup(host string, port string) {
	server.roomnameToRoom = make(map[string]*Chatroom)
	server.roomnameToRoomLock = &sync.Mutex{}

	log.Printf("Setting up server. Will listen via TCP on localhost port %s", port)
	server.networkListener, _ = net.Listen("tcp", ":"+port)
}

// Determines whether or not to let requestingClient either a) Make a new chatroom or b) join an existing chatroom
func (server *Server) resolveChatroomReq(requestingClient *ClientConnection, chatroomReq csprotocol.ChatroomReq) error {
	server.roomnameToRoomLock.Lock()
	defer server.roomnameToRoomLock.Unlock()

	// Make a new empty chatroom if needed
	if chatroomReq.IsNewChatroom {
		_, ok := server.roomnameToRoom[chatroomReq.ChatroomID]
		if ok {
			return errors.New("trying to create a chatroom that already exists")
		}
		newRoom := NewChatroom(chatroomReq.ChatroomID, chatroomReq.ChatroomPassword)
		server.roomnameToRoom[chatroomReq.ChatroomID] = newRoom
	}

	// Add the requestingClient to the existing chatroom
	chatroom, ok := server.roomnameToRoom[chatroomReq.ChatroomID]
	if !ok {
		return errors.New("trying to join a chatroom doesn't exist")
	}
	chatroom.addClient(requestingClient)
	return nil
}

// Broadcasts a message to everyone in the requestingClient's chatroom.
func (server *Server) resolveMessageBroadcastReq(requestingClient *ClientConnection, msgBcstRq csprotocol.MessageBroadcastReq) error {
	if msgBcstRq.LogOut {
		return errors.New("logging out a client")
	}

	server.roomnameToRoomLock.Lock()
	chatroom, _ := server.roomnameToRoom[requestingClient.Roomname]
	server.roomnameToRoomLock.Unlock()

	msgNtf := csprotocol.MessageNotification{Message: msgBcstRq.Message, ClientID: requestingClient.ID}
	chatroom.broadcastToAllExcept(requestingClient, msgNtf)
	return nil
}

func (server *Server) closeConnectionWithClient(closingClient *ClientConnection) {
	server.roomnameToRoomLock.Lock()
	chatroom, _ := server.roomnameToRoom[closingClient.Roomname]
	server.roomnameToRoomLock.Unlock()

	chatroom.removeClient(closingClient)
	closingClient.Connection.Close()
}

// Determines client intent then forever listens for messages (i.e. Message Broadcast Requests) from the client and forwards them.
func (server *Server) handleClient(client *ClientConnection) {
	// Determine the client's identity.
	client.resolveIdentityReq()

	// Determine the chatroom the client wants to join.
	client.resolveRoomReq(
		func(rq *ClientConnection, crq csprotocol.ChatroomReq) error {
			return server.resolveChatroomReq(rq, crq)
		})

	log.Printf("%s has joined %s", client.ID, client.Roomname)

	for {
		// Wait for a Message Broadcast request from the client.
		// Then, broadcast the message in the chatroom.
		err := client.resolveMessageBroadcastReq(
			func(rq *ClientConnection, mrq csprotocol.MessageBroadcastReq) error {
				return server.resolveMessageBroadcastReq(rq, mrq)
			})
		if err != nil {
			server.closeConnectionWithClient(client)
			break
		}
	}
	fmt.Printf("User %s has logged out", client.ID)
}

// Forever listen for new connections.
func (server *Server) registerNewConnections() {
	log.Printf("Waiting for new connections")
	for {
		newConnection, _ := server.networkListener.Accept()
		client := NewClientConnection(newConnection)
		go server.handleClient(client)
	}
}

func (server *Server) terminate() {
	log.Println("\nGoodbye!")
}

// RunLoop spawns a go thread that listens for new client connections.
func (server *Server) RunLoop(quit chan os.Signal) {
	go server.registerNewConnections()
	<-quit
	server.terminate()
}
