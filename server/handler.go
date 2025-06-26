package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Command string `json:"command"`
	Topic   string `json:"topic"`
	Message string `json:"message,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // I wil refactor this later, (CodingTrain reference)
}

func HandleWS(broker *Broker, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrader error:", err)
		return
	}

	client := &Client{
		conn:      conn,
		suscribed: make(map[string]bool),
		send:      make(chan OutgoingMessage, 256),
	}

	go client.writePump()
	log.Println("New Client", r.RemoteAddr)

	defer func() {
		broker.unsubscribeAll <- client
		close(client.send)
		conn.Close()
		log.Println("Client disconected", r.RemoteAddr)
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		switch msg.Command {
		case "SUB":

			if msg.Topic == "" {
				client.send <- OutgoingMessage{
					Topic:   "SYSTEM",
					Message: "Missing Topic",
				}
				continue
			}

			broker.subscribeChan <- subscriptionRequest{
				client: client,
				topic:  msg.Topic,
			}
			client.suscribed[msg.Topic] = true
			client.send <- OutgoingMessage{
				Topic:   "SYSTEM",
				Message: fmt.Sprintf("Subscribed to topic: %s", msg.Topic),
			}
			log.Println(r.RemoteAddr, "Subscribe to -> ", msg.Topic)

		case "UNSUB":
			if msg.Topic == "" {
				client.send <- OutgoingMessage{
					Topic:   "SYSTEM",
					Message: "Missing Topic",
				}
				continue
			}

			broker.unsubscribeChan <- unsubscriptionRequest{
				client: client,
				topic:  msg.Topic,
			}
			delete(client.suscribed, msg.Topic)
			client.send <- OutgoingMessage{
				Topic:   "SYSTEM",
				Message: fmt.Sprintf("Unsubscribed from: %s", msg.Topic),
			}

			log.Println(r.RemoteAddr, "client unsubscribed from:", msg.Topic)

		case "PUB":
			if msg.Topic == "" || msg.Message == "" {
				client.send <- OutgoingMessage{
					Topic:   "System",
					Message: "Missing Topic or Message",
				}
				continue
			}
			broker.publishChan <- publishRequest{
				topic:   msg.Topic,
				message: msg.Message,
			}
			log.Println(r.RemoteAddr, "Message to -> ", msg.Topic, ":", msg.Message)
		default:
			client.send <- OutgoingMessage{
				Topic:   "SYSTEM",
				Message: "Please. Use SUB, UNSUB, or PUB",
			}
		}
	}
}
