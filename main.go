package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	suscribed map[string]bool
}

type Broker struct {
	suscribers map[string]map[*Client]bool
}

type Message struct {
	Command string `json:"command"`
	Topic   string `json:"topic"`
	Message string `json:"message,omitempty"`
}

func newBroker() *Broker {
	return &Broker{
		suscribers: make(map[string]map[*Client]bool),
	}
}

func (b *Broker) Suscribe(topic string, client *Client) {
	if b.suscribers[topic] == nil {
		b.suscribers[topic] = make(map[*Client]bool)
	}
	b.suscribers[topic][client] = true
}

func (b *Broker) Unsuscribe(topic string, client *Client) {
	if subs, ok := b.suscribers[topic]; ok {
		delete(subs, client)
	}
}

func (b *Broker) Publish(topic string, msg string) {
	for client := range b.suscribers[topic] {
		client.SendMessage(topic, msg)
	}
}

func (c *Client) SendMessage(topic, msg string) {
	response := map[string]string{
		"topic":   topic,
		"message": msg,
	}
	c.conn.WriteJSON(response)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // I wil refactor this later, (CodingTrain reference)
}

func handleWS(broker *Broker, w http.ResponseWriter, r *http.Request) {
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
				broker.Suscribe(msg.Topic, client)
				client.suscribed[msg.Topic] = true
				client.SendMessage(" \t >>Server", fmt.Sprintf("Suscribed to %s", msg.Topic))
			}

		case "UNSUB":
			if msg.Topic != "" {
				broker.Unsuscribe(msg.Topic, client)
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

func main() {
	broker := newBroker()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(broker, w, r)
	})

	fmt.Println("Server runnin in :9876/ws")
	log.Fatal(http.ListenAndServe(":9876", nil))
}
