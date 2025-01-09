package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func New(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	router.POST("/users", server.CreateUser)
	router.GET("/users", server.GetUser)
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts", server.ListAccounts)
	router.GET("/accounts/:id", server.GetAccount)
	router.PATCH("/accounts/:id", server.UpdateAccount)
	router.DELETE("/accounts/:id", server.DeleteAccount)
	router.POST("/transfers", server.CreateTransfer)
	router.GET("/transfers/:id", server.GetTransfer)
	router.GET("/transfers", server.ListTransfers)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
