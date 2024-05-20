package socket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/virajBaswana/go_App/utils"
)

var SocketConnections = make(map[*websocket.Conn]map[string]bool)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SocketHandler struct {
	// ws *websocket.Conn
	// database *sqlx.DB
}

func NewSocketConnection(w http.ResponseWriter, r *http.Request) {

	fmt.Println("nwnww")
	Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	fmt.Println("nwnww")
	fmt.Println(r)
	user_id := utils.ExtractClaimsFromRequest(r.Context())

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error")
		fmt.Println(err.Error())
	}
	val := map[string]bool{
		user_id: true,
	}
	SocketConnections[conn] = val
	// fmt.Printf("new connection for user :::---> %v", SocketConnections[user_id])
	fmt.Println(user_id)
	for k, v := range SocketConnections {
		fmt.Println(k)
		fmt.Println(v)
	}
	reader(conn)
}
func reader(ws *websocket.Conn) {
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			// fmt.Println("here")
			// log.Println(err.Error())
			break
		}
		log.Println(string(p))
		if err = ws.WriteMessage(messageType, append([]byte("got ya"), p...)); err != nil {
			fmt.Println("here also")
			log.Println(err.Error())
			break
		}
	}

	ws.Close()
	delete(SocketConnections, ws)
	fmt.Println("post deleteion socket map")
	for k, v := range SocketConnections {
		fmt.Println(k)
		fmt.Println(v)
	}
}
