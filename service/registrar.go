package service

import (
	"github.com/gorilla/websocket"
	"log"
)

type registrar interface {
	register(topic string, options interface{})
}

type websocketRegistrar struct {
	connection *websocket.Conn
}

type AccountTrigger = struct {}

func newWebsocketRegistrar(conn *websocket.Conn) *websocketRegistrar {
	// Register address watching
	watchConfirmations := message{
		Action:  "subscribe",
		Topic:   "confirmation",
		Ack:     false,
		Options: options{
			Accounts: []string{},
		},
	}
	err := conn.WriteJSON(watchConfirmations)
	if err != nil {
		log.Print(err)
	}

	return &websocketRegistrar{
		connection: conn,
	}
}

func (r *websocketRegistrar) register(topic string, options interface{}) {
	r.connection.WriteJSON(message{
		Action:  "update",
		Topic:   topic,
		Ack:     false,
		Options: options,
	})
}
