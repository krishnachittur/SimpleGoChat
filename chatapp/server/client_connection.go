package Server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"log"
	"../csprotocol"
)

type ClientConnection struct {
	ID         string
	Roomname string

	Connection net.Conn
	ConnectionReader     *bufio.Reader
}

func NewClientConnection(newConn net.Conn) *ClientConnection {
	return &ClientConnection {
		Connection: newConn,
		ConnectionReader: bufio.NewReader(newConn),
	}
}

func (cc *ClientConnection) sendMessageNotification(mn csprotocol.MessageNotification) {
	marshalledReq, _ := json.Marshal(mn)
	fmt.Fprintf(cc.Connection, string(marshalledReq)+"\n")
}

// call this method to let the client know the status of the previous request
func (cc *ClientConnection) sendRequestStatus(success bool) {
	reqStatus := csprotocol.RequestStatus{PreviousRequestSucceeded: success}
	marshalledReq, _ := json.Marshal(reqStatus)
	fmt.Fprintf(cc.Connection, string(marshalledReq)+"\n")
}

func (cc *ClientConnection) resolveIdentityReq() {
	identityReq := csprotocol.ClientIdentityReq{}

	// block until a client identity request has been received
	data, _ := cc.ConnectionReader.ReadBytes('\n')
	log.Printf(string(data))
	json.Unmarshal(data, &identityReq)

	// set the client ID
	cc.ID = identityReq.RequestedID
	log.Printf("Setting client ID to %s", cc.ID)

	// acknowledge satisfaction
	cc.sendRequestStatus(true)
}

type chatroomRequestResolver func(*ClientConnection, csprotocol.ChatroomReq) error

func (cc *ClientConnection) resolveRoomReq(resolver chatroomRequestResolver) {
	chatroomReq := csprotocol.ChatroomReq{}

	// block until a ChatroomReq has been received
	data, _ := cc.ConnectionReader.ReadBytes('\n')
	log.Println(string(data))
	json.Unmarshal(data, &chatroomReq)

	// resolve chatroom
	error := resolver(cc, chatroomReq)
	reqSatisfied := error == nil
	fmt.Printf("Resolver returned %t", reqSatisfied)
	if reqSatisfied {
		cc.Roomname = chatroomReq.ChatroomID
		log.Printf("Setting client roomname to %s", cc.Roomname)
	} else {
		log.Printf(error.Error())
	}

	// acknowledge satisfaction
	cc.sendRequestStatus(reqSatisfied)
}

type broadcastRequestResolver func(*ClientConnection, csprotocol.MessageBroadcastReq) error

func (cc *ClientConnection) resolveMessageBroadcastReq(resolver broadcastRequestResolver) {
	broadcastReq := csprotocol.MessageBroadcastReq{}

	// block until a request has been received
	data, _ := cc.ConnectionReader.ReadBytes('\n')
	json.Unmarshal(data, &broadcastReq)

	// resolve broadcast req
	err := resolver(cc, broadcastReq)
	reqSatisfied := err == nil

	// acknowledge satisfaction
	cc.sendRequestStatus(reqSatisfied)
}
