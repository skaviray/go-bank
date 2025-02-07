package api

import (
	"fmt"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func New(store db.Store, config utils.Config) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("unable to create token maker: %c", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: maker,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}
func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.CreateUser)
	router.GET("/users", server.GetUser)
	router.POST("/users/login", server.LoginUser)
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts", server.ListAccounts)
	router.GET("/accounts/:id", server.GetAccount)
	router.PATCH("/accounts/:id", server.UpdateAccount)
	router.DELETE("/accounts/:id", server.DeleteAccount)
	router.POST("/transfers", server.CreateTransfer)
	router.GET("/transfers/:id", server.GetTransfer)
	router.GET("/transfers", server.ListTransfers)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
