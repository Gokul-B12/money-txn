package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Gokul-B12/money-txn/token"
	"github.com/gin-gonic/gin"
)

// declaring defalut values
const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"                //we are using bearer authorization type only
	authorizationPayloadKey = "authorization_payload" //this is new key pair will be stored in gin.contxt
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) { //using anonymous funciton
		authorizationHeaderValue := ctx.GetHeader(authorizationHeaderKey) //GetHeader func returns value of the header key

		if len(authorizationHeaderValue) == 0 {
			err := errors.New("authorization is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		//splitting the string and stores it in slice
		fields := strings.Fields(authorizationHeaderValue)

		//checking the length of fileds.. it must have 2 (i.e, at[0] auth type and [1] access token)
		if len(fields) < 2 {
			err := errors.New("invalid autorization format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		//checking the auth type
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New("invalid autorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]

		//verifying the access token using VerifyToken func from token package,
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		//mapping value of payload to authorizationPayloadKey in gin context
		ctx.Set(authorizationPayloadKey, payload)
		// forwaring this payload to next handler func
		ctx.Next()

	}

}
