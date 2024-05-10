package connection

import (
	"database/sql"
	"time"
)

type Connection struct {
	Id              int       `json:"id"`
	TargetId        int       `json:"target_id"`
	InitiatorId     int       `json:"initiator_id"`
	Is_Reciprocated bool      `json:"is_reciprocated"`
	Created_At      time.Time `json:"created_at"`
	Updated_At      time.Time `json:"updated_at"`
}

type ConnectionModel struct {
	db *sql.DB
}

func (c *ConnectionModel) CreateConnection(body *Connection) (*Connection, error) {
	sqlStmt := `
		INSERT INTO connections (target_id , initiator_id , is_reciprocated)
		VALUES ($1, $2 , $3) RETURNING *;
		`
	row := c.db.QueryRow(sqlStmt, body.TargetId, body.InitiatorId, body.Is_Reciprocated)
	connection := &Connection{}
	if err := row.Scan(&connection.Id, &connection.TargetId, &connection.InitiatorId, &connection.Is_Reciprocated, &connection.Created_At, &connection.Updated_At); err != nil {
		return nil, err
	}
	return connection, nil

}
