package connection

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"viraj_golang/utils"
)

type ConnectionHandler struct {
	router *http.ServeMux
	model  *ConnectionModel
}

func InitConnectionService(mux *http.ServeMux, database *sql.DB) {
	connectionHandler := &ConnectionHandler{
		router: mux,
		model:  &ConnectionModel{db: database},
	}

	connectionHandler.RegisterConnectionRoutes()
}

func (c *ConnectionHandler) RegisterConnectionRoutes() {

	c.router.HandleFunc("POST /connection/create", c.CreateConnection)
}

func (c *ConnectionHandler) CreateConnection(w http.ResponseWriter, r *http.Request) {
	// if err := r.ParseMultipartForm(5242880); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// body := r.MultipartForm

	// var obj = map[string]string{}
	// var obj2 user.User = user.User{}
	connection := &Connection{}
	dec := json.NewDecoder(r.Body)
	// dec.DisallowUnknownFields()
	err := dec.Decode(connection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conn, err := c.model.CreateConnection(connection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SuccessfullyCreated(w, http.StatusCreated, "Created New Connection", map[string]any{"connection": conn})
}
