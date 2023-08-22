package middleware

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
)

var (
	ErrAuthHeaderEmpty           = errors.New("authorization header is not provided")
	ErrAuthHeaderInvalidFormat   = errors.New("invalid authorization header format")
	ErrAuthHeaderUnsupportedType = errors.New("unsupported authorization type")
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

func extractTokenFromHeader(input string) (string, error) {
	if len(input) == 0 {
		return "", ErrAuthHeaderEmpty
	}

	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "", ErrAuthHeaderInvalidFormat
	}

	authorizationType := strings.ToLower(parts[0])
	if authorizationType != authorizationTypeBearer {
		return "", ErrAuthHeaderUnsupportedType
	}

	return parts[1], nil
}

func RequireAuth(repo UserRepository, tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken, err := extractTokenFromHeader(ctx.GetHeader(authorizationHeaderKey))
		if err != nil {
			core.ErrorResponse(ctx, core.ErrUnauthorized.WithError(err.Error()))
			ctx.Abort()
			return
		}

		payload, err := tokenMaker.VerifyToken(authToken)
		if err != nil {
			core.ErrorResponse(ctx, core.ErrUnauthorized.WithError(err.Error()).WithDebug(err.Error()))
			ctx.Abort()
			return
		}

		user, err := repo.GetByUsername(ctx.Request.Context(), payload.UID)
		if err != nil {
			core.ErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		requester := core.NewRequester(strconv.FormatInt(user.ID, 10), user.Username)
		ctx.Set(core.KeyRequester, requester)
		ctx.Next()
	}
}
