package business

import (
	"context"

	"github.com/hoangtk0100/stauvise/pkg/model"
)

type Business interface {
	Auth() AuthBusiness
	User() UserBusiness
}

type AuthBusiness interface {
	Login(ctx context.Context, data *model.LoginParams) (interface{}, error)
}

type UserBusiness interface {
	Register(ctx context.Context, data *model.CreateUserParams) (interface{}, error)
	GetProfile(ctx context.Context) (interface{}, error)
}