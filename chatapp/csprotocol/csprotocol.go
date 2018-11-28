package csprotocol

// ClientIdentityReq is the structure that will be used by the client when requesting for Client IDs
type ClientIdentityReq struct {
	RequestedID string `json:"requested_id"`
}

// ChatroomReq is the structure that will be used by the client when requesting to join or create a chatroom
type ChatroomReq struct {
	ChatroomID       string `json:"chatroom_id"`
	ChatroomPassword string `json:"chatroom_password"`
	IsNewChatroom    bool   `json:"is_new_chatroom"`
}

// MessageBroadcastReq is the structure that will be used by the client when requesting for a message to be broadcasted in the chatroom
type MessageBroadcastReq struct {
	Message string `json:"message"`
	LogOut  bool   `json:"log_out"`
}

// MessageNotification is the structure that will be used by the server to notify a client of a new message in the chatroom
type MessageNotification struct {
	Message  string `json:"message"`
	ClientID string `json:"client_id"`
}

// RequestStatus is the structure that will be used by the server to indicate whether the previous client request succeded or not
type RequestStatus struct {
	PreviousRequestSucceeded bool `json:"prev_req_succeded"`
}
