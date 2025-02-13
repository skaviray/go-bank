package api

import (
	"errors"
	"net/http"
	"simple-bank/token"
	"simple-bank/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationType       = "header"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		authorizationtype := strings.ToLower(fields[0])

		if authorizationtype != authorizationTypeBearer {
			err := errors.New("authorization type is not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
