package business

import (
	"context"
	"errors"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/app-context/util"
	"github.com/hoangtk0100/stauvise/pkg/model"
	"github.com/hoangtk0100/stauvise/pkg/repository"
	"github.com/hoangtk0100/stauvise/pkg/validator"
)

var (
	ErrUserExisted    = errors.New("user existed")
	ErrCannotRegister = errors.New("cannot register")
	ErrCannotGetUser  = errors.New("cannot get user")
)

type userBusiness struct {
	repo repository.Repository
}

func NewUserBusiness(repo repository.Repository) *userBusiness {
	return &userBusiness{repo: repo}
}

func newUserResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         *user.CreatedAt,
	}
}

func (b *userBusiness) Register(ctx context.Context, data *model.CreateUserParams) (*model.UserResponse, error) {
	if err := validator.ValidatePassword(data.Password); err != nil {
		return nil, core.ErrBadRequest.
			WithError(err.Error()).
			WithDebug(err.Error())
	}

	existedUser, err := b.repo.User().GetByUsername(ctx, data.Username)
	if err == nil && existedUser != nil {
		return nil, core.ErrConflict.WithError(ErrUserExisted.Error())
	}

	salt, err := util.RandomString(8)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrCannotRegister.Error()).
			WithDebug(err.Error())
	}

	hashedPassword, err := util.HashPassword("", data.Password, salt)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrCannotRegister.Error()).
			WithDebug(err.Error())
	}

	params := &model.User{
		Username:       data.Username,
		HashedPassword: hashedPassword,
		Salt:           salt,
		Email:          data.Email,
		FullName:       data.FullName,
	}

	user, err := b.repo.User().Create(ctx, params)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrCannotRegister.Error()).
			WithDebug(err.Error())
	}

	return newUserResponse(user), nil
}

func (b *userBusiness) GetProfile(ctx context.Context) (*model.UserResponse, error) {
	requester := core.GetRequester(ctx)

	user, err := b.repo.User().GetByUsername(ctx, requester.GetUID())
	if err != nil {
		return nil, core.ErrUnauthorized.
			WithError(ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	return newUserResponse(user), nil
}
