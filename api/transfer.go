package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/utils"

	"github.com/gin-gonic/gin"
)

type CreateTransferParams struct {
	FromAccount int64  `json:"from_account" binding:"required"`
	ToAccount   int64  `json:"to_account" binding:"required"`
	Amount      int64  `json:"amount" binding:"required,min=1"`
	Currency    string `json:"currency" binding:"required,currency"`
}

func (server *Server) validateCurrency(ctx *gin.Context, accountid int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return account, false
		} else {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return account, false
		}
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] curency mismatch: %s vs %s", accountid, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return account, false
	}
	return account, true
}
func (server *Server) CreateTransfer(ctx *gin.Context) {
	var transfer CreateTransferParams
	if err := ctx.ShouldBindJSON(&transfer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	fromAccount, valid := server.validateCurrency(ctx, transfer.FromAccount, transfer.Currency)
	if !valid {
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(token.Payload)
	if payload.Username != fromAccount.Owner {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}
	_, valid = server.validateCurrency(ctx, transfer.ToAccount, transfer.Currency)
	if !valid {
		return
	}
	args := db.TransferTXParams{
		FromAccount: transfer.FromAccount,
		ToAccount:   transfer.ToAccount,
		Amount:      transfer.Amount,
	}
	_, err := server.store.TransferTX(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)

}

type GetTransferParams struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetTransfer(ctx *gin.Context) {
	var transfer GetTransferParams
	if err := ctx.ShouldBindUri(&transfer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	transferResponse, err := server.store.GetTransfer(ctx, transfer.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
	}
	ctx.JSON(http.StatusOK, transferResponse)

}

type ListTransferParams struct {
	PageId      int32 `form:"page_id" binding:"required,min=1"`
	PageSize    int32 `form:"page_size" binding:"required,min=5,max=10"`
	FromAccount int64 `form:"from_account"`
	ToAccount   int64 `form:"to_account"`
}

func (sever *Server) ListTransfers(ctx *gin.Context) {
	var transfer ListTransferParams
	if err := ctx.ShouldBindQuery(&transfer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	args := db.ListTransfersParams{
		Limit:       transfer.PageId,
		Offset:      transfer.PageSize,
		FromAccount: transfer.FromAccount,
		ToAccount:   transfer.ToAccount,
	}
	transferResponse, err := sever.store.ListTransfers(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transferResponse)
}
