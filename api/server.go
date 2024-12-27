package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func New(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts", server.ListAccounts)
	router.GET("/accounts/:id", server.GetAccount)
	router.PATCH("/accounts/:id", server.UpdateAccount)
	router.DELETE("/accounts/:id", server.DeleteAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
