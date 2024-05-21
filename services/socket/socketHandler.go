package socket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/virajBaswana/go_App/utils"
)

var SocketConnections = make(map[*websocket.Conn]string)

var UserSocketMap = make(map[string]map[*websocket.Conn]bool)

type SendMessage struct {
	Sender_id   int    `json:"sender_id"`
	Receiver_id int    `json:"receiver_id"`
	Message     string `json:"message"`
}

var MessageChannel = make(chan SendMessage)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewSocketConnection(w http.ResponseWriter, r *http.Request) {

	Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	user_id := utils.ExtractClaimsFromRequest(r.Context())

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error")
		fmt.Println(err.Error())
	}

	SocketConnections[conn] = user_id
	UserSocketMap[user_id] = map[*websocket.Conn]bool{conn: true}
	// fmt.Printf("new connection for user :::---> %v", SocketConnections[user_id])
	// fmt.Println(user_id)
	for _, v := range SocketConnections {
		// fmt.Println(k)
		fmt.Println("client %v", v)
	}
	go reader(conn, MessageChannel)
	go writer(MessageChannel)
}
func reader(ws *websocket.Conn, mc chan SendMessage) {
	for {
		message := &SendMessage{}
		err := ws.ReadJSON(message)
		if err != nil {
			ws.Close()
			user_id := SocketConnections[ws]
			delete(SocketConnections, ws)
			delete(UserSocketMap[user_id], ws)
			fmt.Println("post deleteion socket map")
			for _, v := range SocketConnections {
				// fmt.Println(k)
				fmt.Println(v)
			}
			return
		}
		log.Println(message)
		mc <- *message
	}

}

func writer(mc chan SendMessage) {
	for {
		message := <-mc
		receiver := message.Receiver_id
		receiverConnectionsmap := UserSocketMap[strconv.Itoa(receiver)]

		for conn := range receiverConnectionsmap {
			if err := conn.WriteJSON(message); err != nil {
				conn.Close()
				user_id := SocketConnections[conn]
				delete(SocketConnections, conn)
				delete(UserSocketMap[user_id], conn)
				// value = false

				fmt.Println("post deleteion socket map")
				for _, v := range SocketConnections {
					// fmt.Println(k)
					fmt.Println(v)
				}
				return
			}

		}

	}
}
