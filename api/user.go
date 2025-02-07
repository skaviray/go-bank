package api

import (
	"database/sql"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUserReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"fullname"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

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
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	res := newUserResponse(user)
	ctx.JSON(http.StatusOK, res)
}

type GetUserParams struct {
	Username string `form:"username" binding:"required"`
}

func (server *Server) GetUser(ctx *gin.Context) {
	var usersReq GetUserParams
	if err := ctx.ShouldBindQuery(&usersReq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, usersReq.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
	}
	res := newUserResponse(user)
	ctx.JSON(http.StatusOK, res)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
	}
	if err = utils.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
		return
	}
	userRes := newUserResponse(user)
	token, err := server.tokenMaker.CreateToken(user.Username, server.config.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	res := loginUserResponse{
		AccessToken: token,
		User:        userRes,
	}
	ctx.JSON(http.StatusOK, res)
}
