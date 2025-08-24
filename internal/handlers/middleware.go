package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/security"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBearer = "Bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker security.TokenMaker, tokenBlacklist security.TokenBlacklist) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is required")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		isBlacklisted, err := tokenBlacklist.IsTokenBlacklisted(ctx, accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if isBlacklisted {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("token is blacklisted")))
			return
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func GetCurrentUserPayload(ctx *gin.Context) (*security.Payload, error) {
	payload, exists := ctx.Get(AuthorizationPayloadKey)
	if !exists {
		return nil, errors.New("authorization payload not found")
	}

	userPayload, ok := payload.(*security.Payload)
	if !ok {
		return nil, errors.New("invalid authorization payload")
	}

	return userPayload, nil
}
