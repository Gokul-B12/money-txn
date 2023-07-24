package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Gokul-B12/money-txn/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// const (
// 	dbdriver = "postgres"
// 	dbsource = "postgresql://root:admin@34.206.16.110:5432/simple_bank?sslmode=disable"
// )

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatalf("Cannot connect to DB: %v\n", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())

}
