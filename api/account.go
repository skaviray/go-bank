package api

import (
	"database/sql"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/utils"

	"github.com/gin-gonic/gin"
)

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR INR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var params CreateAccountParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	args := db.CreateAccountParams{
		Owner:    params.Owner,
		Currency: params.Currency,
	}
	account, err := server.store.CreateAccount(ctx, args)
	if err != nil {
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
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}

	}
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
	args := db.ListAccountsParams{
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