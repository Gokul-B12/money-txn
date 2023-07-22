package main

import (
	"database/sql"
	"log"

	"github.com/Gokul-B12/money-txn/api"
	db "github.com/Gokul-B12/money-txn/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbdriver      = "postgres"
	dbsource      = "postgresql://root:admin@34.206.16.110:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(dbdriver, dbsource)

	if err != nil {
		log.Fatalf("Cannot connect to DB: %v\n", err)

	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalf("Cannot connect to server: %v\n", err)

	}
}
