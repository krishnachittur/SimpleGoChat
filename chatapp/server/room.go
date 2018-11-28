package Server

import (
	"sync"
	"../csprotocol"
)

type Chatroom struct {
	ID string
	password string

	clients []*ClientConnection
	clientsLock *sync.Mutex
}

func NewChatroom(ID string, password string) *Chatroom {
	return &Chatroom{
		ID: ID, password: password,
		clients: make([]*ClientConnection, 0),
		clientsLock: &sync.Mutex{},
	}
}

func (room *Chatroom) broadcastToAllExcept(excludedClient *ClientConnection, msgNtf csprotocol.MessageNotification) {
	room.clientsLock.Lock()
	defer room.clientsLock.Unlock()

	for _, client := range room.clients {
		if client.ID == excludedClient.ID {
			continue
		}
		client.sendMessageNotification(msgNtf)
	}
}

func (room *Chatroom) addClient(newClient *ClientConnection) {
	room.clientsLock.Lock()
	defer room.clientsLock.Unlock()

	room.clients = append(room.clients, newClient)
}

