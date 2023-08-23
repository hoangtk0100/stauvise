package business

import (
	"context"
	"mime/multipart"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

type Business interface {
	Auth() AuthBusiness
	User() UserBusiness
	Category() CategoryBusiness
	Video() VideoBusiness
}

type AuthBusiness interface {
	Login(ctx context.Context, data *model.LoginParams) (interface{}, error)
}

type UserBusiness interface {
	Register(ctx context.Context, data *model.CreateUserParams) (interface{}, error)
	GetProfile(ctx context.Context) (interface{}, error)
}

type CategoryBusiness interface {
	Create(ctx context.Context, data *model.CreateCategoryParams) (*model.Category, error)
	GetAll(ctx context.Context) ([]model.Category, error)
}

type VideoBusiness interface {
	Create(ctx context.Context, data *model.CreateVideoParams, fileHeader *multipart.FileHeader) (uint64, error)
	GetByID(ctx context.Context, id uint64) (interface{}, error)
	List(ctx context.Context, filter *model.VideoFilter, paging *core.Paging) ([]model.Video, error)
	GetSegments(ctx context.Context, id uint64, paging *core.Paging) (interface{}, error)
}
