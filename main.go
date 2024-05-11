package main

import (
	"fmt"
	"log"
	"net/http"
	"viraj_golang/db"
	"viraj_golang/middlewares"
	"viraj_golang/services/auth"
	"viraj_golang/services/connection"
	"viraj_golang/services/user"
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
	//registering all the routes and paths
	//sub routers
	authRouter := auth.InitAuthRoutes(db)
	userRouter := user.InitUserRoutes(db)
	connectionRouter := connection.InitConnectionService(db)

	// global middleware stack on all routes , main router
	middlewareStack := middlewares.CreateMiddlewareStack(
		middlewares.RecoverPanic,
		middlewares.RequestLogger,
		middlewares.SecureHeaders,
	)

	//integrate all sub routers into the main router
	mux.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	mux.Handle("/user/", http.StripPrefix("/user", middlewares.CheckAuth(userRouter)))
	mux.Handle("/connection/", http.StripPrefix("/connection", middlewares.CheckAuth(connectionRouter)))

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewareStack(mux),
	}
	log.Fatal(server.ListenAndServe())
	fmt.Println("server online")
}
