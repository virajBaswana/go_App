package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func EstablishandVerifyDBConnection() (*sql.DB, error) {
	const connectionString = "user=postgres dbname=go_practice password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		// log.Fatal("DB connection failed")
		return nil, err
	}
	if err = db.Ping(); err != nil {
		// log.Fatal("DB Ping Failed")
		return nil, err
	}
	fmt.Println("db connected")
	return db, nil

}
