package server

import "github.com/gorilla/websocket"

type Client struct {
	conn      *websocket.Conn
	suscribed map[string]bool
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:      conn,
		suscribed: make(map[string]bool),
	}
}

func (c *Client) SendMessage(topic, msg string) {
	response := map[string]string{
		"topic":   topic,
		"message": msg,
	}
	c.conn.WriteJSON(response)
}
