package business

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/stauvise/pkg/model"
	"github.com/hoangtk0100/stauvise/pkg/repository"
	"github.com/hoangtk0100/stauvise/pkg/validator"
)

var (
	ErrEmailOrPasswordInvalid = errors.New("email or password invalid")
)

type authBusiness struct {
	repo       repository.Repository
	tokenMaker token.TokenMaker
}

func NewAuthBusiness(repo repository.Repository, tokenMaker token.TokenMaker) *authBusiness {
	return &authBusiness{
		repo:       repo,
		tokenMaker: tokenMaker,
	}
}

type loginResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (a *authBusiness) Login(ctx context.Context, data *model.LoginParams) (interface{}, error) {
	user, err := a.repo.User().GetByUsername(ctx, data.Username)
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return nil, core.ErrNotFound.
				WithError(ErrEmailOrPasswordInvalid.Error())
		}

		return nil, core.ErrInternalServerError.
			WithError(ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	if err := validator.ValidatePassword(data.Password); err != nil {
		return nil, core.ErrBadRequest.
			WithError(err.Error()).
			WithDebug(err.Error())
	}

	err = util.CheckPassword("", user.HashedPassword, data.Password, user.Salt)
	if err != nil {
		return nil, core.ErrUnauthorized.
			WithError(ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	accessToken, accessPayload, err := a.tokenMaker.CreateToken(token.AccessToken, user.Username)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(token.RefreshToken, user.Username)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, core.ErrInternalServerError.
			WithError(ErrEmailOrPasswordInvalid.Error())
	}

	session, err := a.repo.Session().Create(ctx, &model.Session{
		ID:           refreshPayload.ID,
		OwnerID:      user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ginCtx.Request.UserAgent(),
		ClientIp:     ginCtx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrEmailOrPasswordInvalid.Error()).
			WithDebug(err.Error())
	}

	rsp := loginResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	return rsp, nil
}
