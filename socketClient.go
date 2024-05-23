package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/virajBaswana/go_App/services/message"
	"github.com/virajBaswana/go_App/utils"
)

type SocketClient struct {
	user_id        int
	socket         *websocket.Conn
	sockethub      *SocketHub
	receiveMessage chan *message.SocketMessage
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *SocketClient) reader() {
	for {
		message := &message.SocketMessage{}
		if err := c.socket.ReadJSON(message); err != nil {
			log.Println(err)
			c.sockethub.removeClient <- c
			c.socket.Close()
			break
		}
		c.sockethub.sendMessage <- message
	}
}
func (c *SocketClient) writer() {
	defer func() {
		delete(c.sockethub.clients, c)
		close(c.receiveMessage)
	}()
	for {
		messagetoSend, _ := <-c.receiveMessage
		if err := c.socket.WriteJSON(messagetoSend); err != nil {
			log.Println(err)
			break
		}
		n := len(c.receiveMessage)
		for i := 0; i < n; i++ {
			if err := c.socket.WriteJSON(messagetoSend); err != nil {
				log.Println(err)
				break
			}
		}

	}
}

func ServeWs(hub *SocketHub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	user := utils.ExtractClaimsFromRequest(r.Context())
	userId, _ := strconv.Atoi(user)
	client := &SocketClient{
		user_id:        userId,
		sockethub:      hub,
		socket:         connection,
		receiveMessage: make(chan *message.SocketMessage, 256),
	}
	client.sockethub.addClient <- client

	go client.reader()
	go client.writer()
}
