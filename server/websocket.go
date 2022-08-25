package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	users = make(map[string]*websocket.Conn)
	chat  = make(chan NewMessage)
)


func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// fmt.Println("p Before", string(p))

		var details NewMessage

		errMarsh := json.Unmarshal(p, &details)

		if errMarsh != nil {
			fmt.Println("Error unmarshalling: ", errMarsh)
		}

		details.messageType = messageType

	
		// print out that message for clarity
		fmt.Println("Msg revieved: ", details)

		// Add To the channel instead of writing the message back
		if _, recieverValid := users[details.RecieverID]; !recieverValid {
			fmt.Println("User sending two isnt active UserID: ", details.RecieverID)
		} else {
			chat <- details
		}
	}
}

func SendMsgs() {
	for {
		select {
		case msg, ok := <-chat:
			if ok {
				// fmt.Println("Attempting to send: ", msg)
					//Add consition to check the user exsists
					if _, valid := users[msg.RecieverID]; valid {
						sendMess := SendMessage{Sender: msg.UserID, Message: msg.Mesg}
						// fmt.Println("SENDING BACK", sendMess)
						res, marshErr := json.Marshal(sendMess)
						if marshErr != nil {
							fmt.Println("Error MArshalling the data before sending: ", marshErr)
						}
						err := users[msg.RecieverID].WriteMessage(msg.messageType, res)
						// fmt.Println(users[msg.RecieverID])
						if err != nil {
							log.Printf("error: %v", err)
								users[msg.RecieverID].Close()
								delete(users, msg.RecieverID)
							return
						}
						fmt.Println("Message sent to: ", msg.RecieverID)

					}
			} else {
				chat = nil
			}
		}
	}
}