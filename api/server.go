package api

import (
	_ "github.com/Gokul-B12/money-txn/db"
)

// this server serves all our HTTP requests for our banking service.
type server struct {
	store
}
