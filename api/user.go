package api

import (
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/utils"

	"github.com/gin-gonic/gin"
)

type CreateUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6	"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var userReq CreateUserReq
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	hashedPass, err := utils.CreateHashedPassword(userReq.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	args := db.CreateUserParams{
		Username:       userReq.Username,
		HashedPassword: hashedPass,
		Email:          userReq.Email,
		FullName:       userReq.FullName,
	}
	user, err := server.store.CreateUser(ctx, args)
	if err != nil {
		utils.ErrorResponse(err)
	}
	ctx.JSON(http.StatusOK, user)
}

func (server *Server) GetUser(ctx *gin.Context) {

}
