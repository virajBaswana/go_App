package connection

import (
	"time"

	"github.com/jmoiron/sqlx"
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
	db *sqlx.DB
}

func (c *ConnectionModel) CreateConnection(body *Connection) (*Connection, error) {
	sqlStmt := `
		INSERT INTO connections (target_id , initiator_id , is_reciprocated)
		VALUES ($1, $2 , $3) RETURNING *;
		`
	row := c.db.QueryRowx(sqlStmt, body.TargetId, body.InitiatorId, body.Is_Reciprocated)
	connection := &Connection{}
	if err := row.StructScan(connection); err != nil {
		return nil, err
	}
	return connection, nil

}
