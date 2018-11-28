package Server

import (
	"sync"
	"../csprotocol"
)

type Chatroom struct {
	ID string
	password string

	clients []*ClientConnection
	clientsLock sync.Mutex
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
