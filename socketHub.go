package main

import "github.com/virajBaswana/go_App/services/message"

type SocketHub struct {
	clients      map[*SocketClient]bool
	sendMessage  chan *message.SocketMessage
	addClient    chan *SocketClient
	removeClient chan *SocketClient
}

func newSocketHub() *SocketHub {
	return &SocketHub{
		sendMessage:  make(chan *message.SocketMessage),
		addClient:    make(chan *SocketClient),
		removeClient: make(chan *SocketClient),
		clients:      make(map[*SocketClient]bool),
	}
}

// func (hub *SocketHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {

// }

func (hub *SocketHub) run() {
	for {
		select {
		case client := <-hub.addClient:
			hub.clients[client] = true
		case client := <-hub.removeClient:
			delete(hub.clients, client)
			close(client.receiveMessage)
		case message := <-hub.sendMessage:
			for client := range hub.clients {
				if client.user_id == message.Receiver_id {
					select {
					case client.receiveMessage <- message:
					default:
						delete(hub.clients, client)
						close(client.receiveMessage)
					}

				}
			}
		}
	}
}
