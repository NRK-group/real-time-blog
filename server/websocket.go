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

		var details NewMessage

		errMarsh := json.Unmarshal(p, &details)

		if errMarsh != nil {
			fmt.Println("Error unmarshalling: ", errMarsh)
			return
		}

		details.messageType = messageType


		// Add To the channel instead of writing the message back
		if _, recieverValid := users[details.RecieverID]; !recieverValid {
			fmt.Println("User sending two isnt active UserID: ", details.RecieverID)
		} else {
			chat <- details
		}
		//Now add the messafe to the database

		
	}
}

func SendMsgs() {
	for {
		select {
		case msg, ok := <-chat:
			if ok {
				// Add consition to check the user exsists
				if _, valid := users[msg.RecieverID]; valid {
					sendMess := SendMessage{Sender: msg.UserID, Message: msg.Mesg, Date: msg.Date}
					res, marshErr := json.Marshal(sendMess)
					if marshErr != nil {
						fmt.Println("Error MArshalling the data before sending: ", marshErr)
						return
					}
					err := users[msg.RecieverID].WriteMessage(msg.messageType, res)
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
