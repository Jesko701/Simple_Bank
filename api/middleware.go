package api

import (
	"errors"
	"fmt"
	"net/http"
	"solo_simple-bank_tutorial/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader_Key  = "authorization"
	AuthorizationType_Bearer = "Bearer"
	Authorization_Payload    = "authorization_payload"
)

func authMiddleware(token token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the header
		authorizationHeader := ctx.GetHeader(AuthorizationHeader_Key)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		// Split the string of the header
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		// Check the bearer type
		authorizationType := fields[0]
		if strings.ToLower(authorizationType) != strings.ToLower(AuthorizationType_Bearer) {
			err := fmt.Errorf("invalid authorization type : %v", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		//Check the validation of the token
		accessToken := fields[1]
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		ctx.Set(Authorization_Payload, payload)
		ctx.Next()
	}
}
