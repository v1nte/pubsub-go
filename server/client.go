package server

import (
	"github.com/gorilla/websocket"
	"github.com/v1nte/pubsub-go/logger"
	"go.uber.org/zap"
)

type OutgoingMessage struct {
	Author  string
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
		name:      "",
		suscribed: make(map[string]bool),
		send:      make(chan OutgoingMessage, 256),
	}
}

func (c *Client) writePump() {
	logger.Log.Info("Initialize writePump for", zap.Any("Client: ", c))
	defer func() {
		logger.Log.Info("Close writePump for", zap.Any("Client: ", c))
	}()

	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			logger.Log.Error("error in c.writePump", zap.Error(err))
		}
	}
}
