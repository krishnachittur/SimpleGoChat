package Client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"../csprotocol"
)

// Client implements the Endpoint interface
type Client struct {
	connection    net.Conn
	networkReader *bufio.Reader
	stdinReader   *bufio.Reader
}

func (client *Client) Setup(host string, port string) {
	client.connection, _ = net.Dial("tcp", host+":"+port)
	client.networkReader = bufio.NewReader(client.connection)
	client.stdinReader = bufio.NewReader(os.Stdin)
	client.getNewClientID()
	client.joinChatroom()
}

func (client *Client) readInput() string {
	inputStr, _ := client.stdinReader.ReadString('\n')
	return strings.TrimSpace(inputStr)
}

func (client *Client) getNewClientID() {
	fmt.Print("Enter the User Name you want: ")
	for {
		// determine parameters
		identityRequest := csprotocol.ClientIdentityReq{}
		identityRequest.RequestedID = client.readInput()

		// request username
		reqMarshalled, err := json.Marshal(identityRequest)
		if err != nil {
			log.Fatal("Couldn't marshal identityRequest")
		}

		fmt.Fprintf(client.connection, string(reqMarshalled)+"\n")

		// check to make sure request succeded
		reqStatus, _ := client.getRequestStatus()
		if !reqStatus.PreviousRequestSucceeded {
			fmt.Print("Couldn't get requested username. Please try another one:")
			continue
		}
		break
	}
}

func (client *Client) joinChatroom() {
	for {
		// determine parameters
		chatroomReq := csprotocol.ChatroomReq{}

		fmt.Print("Press 'y' if you want to make a new chatroom and 'n' if you want to join an existing one:")
		yOrN := client.readInput()
		if strings.ToLower(yOrN) == "y" {
			chatroomReq.IsNewChatroom = true
		} else if strings.ToLower(yOrN) == "n" {
			chatroomReq.IsNewChatroom = false
		} else {
			fmt.Println("You didn't type 'y' or 'n'. Retry.")
			continue
		}

		fmt.Print("Enter the name of the chatroom you want to join or create:")
		chatroomReq.ChatroomID = client.readInput()

		fmt.Print("Enter the password of the chatroom you want to join or create:")
		chatroomReq.ChatroomPassword = client.readInput()

		// request chatroom
		reqMarshalled, _ := json.Marshal(chatroomReq)
		fmt.Fprintf(client.connection, string(reqMarshalled)+"\n")

		// check to make sure request succeded
		reqStatus, _ := client.getRequestStatus()
		if !reqStatus.PreviousRequestSucceeded {
			fmt.Println("Couldn't make or join requested chatroom.")
			fmt.Println("If you're trying to make a new chatroom, the chatroom name is likely already used. Please try a new name.")
			fmt.Println("If you're trying to join an existing chatroom, you either misspelled the chatroom name or the password's incorrect.")
			fmt.Println("Starting retry procedure.")
			continue
		}
		break
	}
}

func (client *Client) getRequestStatus() (reqStatus csprotocol.RequestStatus, err error) {
	data, err := client.networkReader.ReadBytes('\n')
	if err != nil {
		log.Fatal("error in the received request status")
	}

	err = json.Unmarshal(data, &reqStatus)
	if err != nil {
		log.Fatal("error when unmarshalling request status")
	}
	return
}

func (client *Client) getMessageNotification() (mn csprotocol.MessageNotification, err error) {
	data, err := client.networkReader.ReadBytes('\n')
	if err != nil {
		log.Fatal("error in the received message")
	}

	err = json.Unmarshal(data, &mn)
	if err != nil {
		log.Fatal("error when unmarshalling received message")
	}
	return
}

func (client *Client) listenLoop() {
	for {
		mn, err := client.getMessageNotification()
		if err != nil {
			continue
		}

		clientID := strings.TrimSpace(mn.ClientID)
		message := strings.TrimSpace(mn.Message)
		if len(message) > 0 {
			fmt.Println(clientID + ": " + message)
		}
	}
}

func (client *Client) sendLoop() {
	for {
		// determine parameters
		messageReq := csprotocol.MessageBroadcastReq{LogOut: false}

		fmt.Print("> ")
		messageReq.Message = client.readInput()

		// request message to be broadcasted
		reqMarshalled, _ := json.Marshal(messageReq)
		fmt.Fprintf(client.connection, string(reqMarshalled)+"\n")
	}
}

func (client *Client) terminate() {
	// determine parameters
	messageReq := csprotocol.MessageBroadcastReq{LogOut: true, Message: "Goodbye!"}

	// request termination of client connection
	reqMarshalled, _ := json.Marshal(messageReq)
	fmt.Fprintf(client.connection, string(reqMarshalled)+"\n")
}

func (client *Client) RunLoop(quit chan os.Signal) {
	go client.listenLoop()
	go client.sendLoop()
	<-quit
	client.terminate()
}
