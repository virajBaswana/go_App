package connection

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Connection struct {
	Id              int       `json:"id"`
	Target_Id       int       `json:"target_id"`
	Initiator_Id    int       `json:"initiator_id"`
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

	row := c.db.QueryRowx(sqlStmt, body.Target_Id, body.Initiator_Id, body.Is_Reciprocated)
	connection := &Connection{}
	if err := row.StructScan(connection); err != nil {
		return nil, err
	}
	return connection, nil
}

func (c *ConnectionModel) InverseConnectionRecord(body *Connection) (*Connection, error) {
	sqlStatement := `SELECT * FROM connections WHERE target_id=$1 AND initiator_id=$2;`
	row := c.db.QueryRowx(sqlStatement, body.Initiator_Id, body.Target_Id)
	conn := &Connection{}
	// var conn *Connection

	if err := row.StructScan(conn); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println(err.Error())
		return nil, err
	}
	return conn, nil
}

func (c *ConnectionModel) UpdateConnection(body *Connection) (*Connection, error) {
	q := `UPDATE connections SET is_reciprocated=$1 WHERE initiator_id=$2 AND target_i=$3 RETURNING *;`
	conn := &Connection{}
	if err := c.db.QueryRowx(q, body.Is_Reciprocated, body.Initiator_Id, body.Target_Id).StructScan(conn); err != nil {
		return nil, err
	}
	return conn, nil
}
