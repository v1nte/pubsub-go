package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type OutgoingMessage struct {
	Topic   string
	Message string
}

type Client struct {
	conn      *websocket.Conn
	suscribed map[string]bool
	send      chan OutgoingMessage
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:      conn,
		suscribed: make(map[string]bool),
		send:      make(chan OutgoingMessage, 256),
	}
}

func (c *Client) SendMessage(topic, msg string) {
	response := map[string]string{
		"topic":   topic,
		"message": msg,
	}
	c.conn.WriteJSON(response)
}

func (c *Client) writePump() {
	for msg := range c.send {
		err := c.conn.WriteJSON(map[string]string{
			"topic":   msg.Topic,
			"message": msg.Message,
		})
		if err != nil {
			log.Println("Some error in c.writePump", err)
		}
	}
}
