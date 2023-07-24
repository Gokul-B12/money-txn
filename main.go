package main

import (
	"database/sql"
	"log"

	"github.com/Gokul-B12/money-txn/api"
	db "github.com/Gokul-B12/money-txn/db/sqlc"
	"github.com/Gokul-B12/money-txn/util"
	_ "github.com/lib/pq"
)

// const (
// 	dbdriver      = "postgres"
// 	dbsource      = "postgresql://root:admin@34.206.16.110:5432/simple_bank?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatalf("Cannot connect to DB: %v\n", err)

	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("Cannot connect to server: %v\n", err)

	}
}
