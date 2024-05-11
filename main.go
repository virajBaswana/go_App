package main

import (
	"fmt"
	"log"
	"net/http"
	"viraj_golang/db"
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

	mux := http.NewServeMux()
	//registering all the routes and paths
	user.InitUserRoutes(mux, db)
	connection.InitConnectionService(mux, db)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
	fmt.Println("server online")
}
