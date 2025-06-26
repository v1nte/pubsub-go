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
	name      string
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

func (c *Client) writePump() {
	log.Println("initialize writePump for", c)
	defer func() {
		log.Println("close writePump for", c)
	}()

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
