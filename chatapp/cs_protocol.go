package main

type ClientIdentityReq struct {
	RequestedID string `json:"requested_id"`
}

type ChatroomReq struct {
	ChatroomID       string `json:"chatroom_id"`
	ChatroomPassword string `json:"chatroom_password"`
	IsNewChatroom    bool   `json:"is_new_chatroom"`
}

type MessageReq struct {
	Message string `json:"message"`
	LogOut  bool   `json:"log_out"`
}

type ErrorRes struct {
	ErrorType int `json:"error"`
}
