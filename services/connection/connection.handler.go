package connection

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/virajBaswana/go_App/middlewares"
	"github.com/virajBaswana/go_App/utils"

	"github.com/jmoiron/sqlx"
)

type ConnectionHandler struct {
	router *http.ServeMux
	model  *ConnectionModel
}

func InitConnectionService(database *sqlx.DB) *http.ServeMux {
	connectionRouter := http.NewServeMux()
	middlewares.CheckAuth(connectionRouter)
	connectionHandler := &ConnectionHandler{
		router: connectionRouter,
		model:  &ConnectionModel{db: database},
	}

	connectionHandler.RegisterConnectionRoutes()
	return connectionRouter
}

func (c *ConnectionHandler) RegisterConnectionRoutes() {

	c.router.HandleFunc("POST /createConnection", c.CreateConnection)
}

func (c *ConnectionHandler) CreateConnection(w http.ResponseWriter, r *http.Request) {

	connection := &Connection{}
	dec := json.NewDecoder(r.Body)
	// dec.DisallowUnknownFields()
	err := dec.Decode(connection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	inverserConnection, err := c.model.InverseConnectionRecord(connection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("aaaaaadfasfas")
		return
	}
	if inverserConnection != nil {
		fmt.Println("aaaaajhbkhjh")
		utils.SuccessfullyFoundOne(w, &utils.JsonResponse{Message: "connection alreafy exists", Code: http.StatusFound, Body: map[string]any{
			"existing_connection": inverserConnection,
		}})

	} else {
		conn, err := c.model.CreateConnection(connection)
		if err != nil {
			fmt.Println("aaaaa")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		utils.SuccessfullyCreated(w, http.StatusCreated, "Created New Connection", map[string]any{"connection": conn})
	}
}
