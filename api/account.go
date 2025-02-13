package api

import (
	"database/sql"
	"errors"

	// "go/token"
	"log"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountParams struct {
	// Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR INR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var params CreateAccountParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.CreateAccountParams{
		Owner:    payload.Username,
		Currency: params.Currency,
	}
	account, err := server.store.CreateAccount(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type GetAccountParams struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var params GetAccountParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, params.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != payload.Username {
		err := errors.New("account doesnt belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	// account = db.Account{}
	ctx.JSON(http.StatusOK, account)
}

type ListAccountParams struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccounts(ctx *gin.Context) {
	var params ListAccountParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  params.PageSize,
		Offset: (params.PageId - 1) * params.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) DeleteAccount(ctx *gin.Context) {

}

func (server *Server) UpdateAccount(ctx *gin.Context) {

}
