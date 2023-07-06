package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	dbdriver = "postgres"
	dbsource = "postgresql://root:admin@34.206.16.110:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbdriver, dbsource)

	if err != nil {
		log.Fatalf("Cannot connect to DB: %v\n", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())

}
