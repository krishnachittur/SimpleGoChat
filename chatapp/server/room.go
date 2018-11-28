package Server

import (
	"sync"
	"../csprotocol"
)

type Chatroom struct {
	ID string
	password string

	clients []ClientConnection
	clientsLock sync.Mutex
}

func (room Chatroom) broadcastToAllExcpet(exclude ClientConnection, msgBcstRq csprotocol.MessageBroadcastReq) {
	room.clientsLock()
	defer room.clientsUnlock()

	for _, client in range room.clients() {
		if client.ID == exclude.ID {
			continue
		}
		client.sendNotification(msgBcstRq)
	}
}
