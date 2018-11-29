package Server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"../csprotocol"
)

type ClientConnection struct {
	ID       string
	Roomname string

	Connection       net.Conn
	ConnectionReader *bufio.Reader
}

func NewClientConnection(newConn net.Conn) *ClientConnection {
	return &ClientConnection{
		Connection:       newConn,
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
	json.Unmarshal(data, &identityReq)

	// set the client ID
	cc.ID = identityReq.RequestedID

	// acknowledge satisfaction
	cc.sendRequestStatus(true)
}

type chatroomRequestResolver func(*ClientConnection, csprotocol.ChatroomReq) error

func (cc *ClientConnection) resolveRoomReq(resolver chatroomRequestResolver) {
	reqSatisfied := false
	for !reqSatisfied {
		chatroomReq := csprotocol.ChatroomReq{}

		// block until a ChatroomReq has been received
		data, _ := cc.ConnectionReader.ReadBytes('\n')
		json.Unmarshal(data, &chatroomReq)

		// resolve chatroom
		error := resolver(cc, chatroomReq)
		reqSatisfied = error == nil
		if reqSatisfied {
			cc.Roomname = chatroomReq.ChatroomID
		}

		// acknowledge satisfaction
		cc.sendRequestStatus(reqSatisfied)
	}
}

type broadcastRequestResolver func(*ClientConnection, csprotocol.MessageBroadcastReq) error

func (cc *ClientConnection) resolveMessageBroadcastReq(resolver broadcastRequestResolver) error {
	broadcastReq := csprotocol.MessageBroadcastReq{}

	// block until a request has been received
	data, _ := cc.ConnectionReader.ReadBytes('\n')
	err := json.Unmarshal(data, &broadcastReq)

	// this error will happen only if the data being sent isn't of type MessageBroadcastReq
	// this will happen only if the client wants to log out
	if err != nil {
		return err
	}

	// resolve broadcast req
	err = resolver(cc, broadcastReq)
	if err != nil {
		log.Println("Error: " + err.Error())
	}
	// TODO: maybe in the future support acknowledgements on messages as well
	// cc.sendRequestStatus(reqSatisfied)
	return nil
}
