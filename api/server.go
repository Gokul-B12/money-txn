package api

import (
	"fmt"

	db "github.com/Gokul-B12/money-txn/db/sqlc"
	"github.com/Gokul-B12/money-txn/token"
	"github.com/Gokul-B12/money-txn/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// this server serves all our HTTP requests for our banking service.
type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %s", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()

	return server, nil

}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.POST("/users/login", server.loginUser)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	//router.PUT("/accounts", server.updateAccount)
	//router.DELETE("/accounts/:id", server.deleteAccount)
	router.POST("/transfers", server.createTransfer)
	router.POST("/users", server.createUser)

	server.router = router
}

func (server *Server) Start(address string) error {

	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}

}
