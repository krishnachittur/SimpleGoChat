package Server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"../csprotocol"
)

type ClientConnection struct {
	ID         string
	Roomname string

	Connection net.Conn
	ConnectionReader     *bufio.Reader
}

func (cc ClientConnection) sendMessageNotfication(mn csprotocol.MessageNotification) {
	marshalledReq, _ := json.Marshal(mn)
	fmt.Fprintf(cc.Connection,, string(marshalledReq)+"\n")
}

// call this method to let the client know the status of the previous request
func (cc ClientConnection) sendRequestStatus(success bool) {
	reqStatus := csprotocol.RequestStatus{PreviousRequestSucceeded: success}
	marshalledReq, _ := json.Marshal(reqStatus)
	fmt.Fprintf(cc.Connection, string(marshalledReq)+"\n")
}

func (cc ClientConnection) resolveIdentityRequest() {
	identityReq := csprotocol.ClientIdentityReq{}

	// block until a client identity request has been received
	data, _ := cc.ConnectionReader.ReadBytes('\n')
	json.Unmarshal(data, &identityReq)

	// set the client ID
	cc.clientID = identityReq.RequestedID

	// acknowledge satisfaction
	cc.sendRequestStatus(true)
}

type chatroomRequestResolver func(ClientConnection, csprotocol.ChatroomReq) error

func (cc ClientConnection) resolveRoomReq(resolver chatroomRequestResolver) {
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

type broadcastRequestResolver func(ClientConnection, csprotocol.MessageBroadcastReq) error

func (cc ClientConnection) resolveMessageBroadcastReq(resolver broadcastRequestResolver) {
	broadcastReq := csprotocol.MessageBroadcastReq{}

	// block until a request has been received
	data, _ := cc.ConnectionReader.ReadBytes('\n')
	json.Unmarshal(data, &broadcastReq)

	// resolve broadcast req
	err := resolver(cc, broadcastRequestResolver)
	reqSatisfied := err == nil

	// acknowledge satisfaction
	cc.sendRequestStatus(reqSatisfied)
}
