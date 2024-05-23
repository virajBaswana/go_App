package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/virajBaswana/go_App/db"
	"github.com/virajBaswana/go_App/middlewares"
	"github.com/virajBaswana/go_App/services/auth"
	"github.com/virajBaswana/go_App/services/connection"

	"github.com/virajBaswana/go_App/services/user"
)

func main() {
	db, err := db.EstablishandVerifyDBConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	// closing db for graceful shutdown
	defer db.Close()

	// main router
	mux := http.NewServeMux()

	//web sockets
	sockethub := newSocketHub()
	go sockethub.run()
	//registering all the routes and paths
	//sub routers
	authRouter := auth.InitRoutes(db)
	userRouter := user.InitRoutes(db)
	connectionRouter := connection.InitRoutes(db)

	// socket := socket.
	// subRouters := []SubRouter{
	// 	authRouter, userRouter, connectionRouter,
	// }

	// global middleware stack on all routes , main router
	middlewareStack := middlewares.CreateMiddlewareStack(
		middlewares.RecoverPanic,
		middlewares.RequestLogger,
		middlewares.SecureHeaders,
	)
	// socketHandlerForMiddlewareChaining :=

	//integrate all sub routers into the main router
	mux.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	mux.Handle("/user/", http.StripPrefix("/user", middlewares.CheckAuth(userRouter)))
	mux.Handle("/connection/", http.StripPrefix("/connection", middlewares.CheckAuth(connectionRouter)))

	mux.Handle("/ws", middlewares.CheckAuth(http.HandlerFunc(socket.NewSocketConnection)))

	// upgrader := websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// }
	// func wshandler(w http.ResponseWriter, r *http.Request) {
	// 	conn, err := upgrader.Upgrade(w, r, nil)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	for {
	// 		messageType, p, err := conn.ReadMessage()
	// 		if err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 		if err := conn.WriteMessage(messageType, p); err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 	}
	// }

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewareStack(mux),
	}
	log.Fatal(server.ListenAndServe())
	fmt.Println("server online")
}
