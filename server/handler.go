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
		log.Println("Upgrader:", err)
	}

	client := &Client{
		conn:      conn,
		suscribed: make(map[string]bool),
	}

	defer func() {
		conn.Close()
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
			if msg.Topic != "" {
				broker.Subscribe(msg.Topic, client)
				client.suscribed[msg.Topic] = true
				client.SendMessage(" \t >>Server", fmt.Sprintf("Suscribed to %s", msg.Topic))
			}

		case "UNSUB":
			if msg.Topic != "" {
				broker.Unsubscribe(msg.Topic, client)
				delete(client.suscribed, msg.Topic)
				client.SendMessage(" \t >>Server", fmt.Sprintf("Unsusbribed to %s", msg.Topic))
			}

		case "PUB":
			if msg.Topic != "" && msg.Message != "" {
				broker.Publish(msg.Topic, msg.Message)
			}
		default:
			client.SendMessage(" \t >>Server", "Use «SUB», «UNSUB», or «PUB»")
		}
	}
}
