package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func EstablishandVerifyDBConnection() (*sqlx.DB, error) {
	const connectionString = "user=postgres dbname=go_practice password=postgres sslmode=disable"
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		// log.Fatal("DB connection failed")
		return nil, err
	}
	if err = db.Ping(); err != nil {
		// log.Fatal("DB Ping Failed")
		return nil, err
	}
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://postgres:postgres@localhost:5432/go_practice?sslmode=disable")
	if err != nil {
		fmt.Println("asdfasdfa")
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("asdfasdfa")
		log.Fatal(err)
	}
	fmt.Println("db connected and migrations done")
	return db, nil

}
