package repository

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

type Repository interface {
	User() UserRepository
	Session() SessionRepository
	Category() CategoryRepository
	Video() VideoRepository
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, data *model.User) (*model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, data *model.Session) (*model.Session, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, data *model.Category) (*model.Category, error)
	GetByName(ctx context.Context, name string) (*model.Category, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]model.Category, error)
	GetAll(ctx context.Context) ([]model.Category, error)
}

type VideoRepository interface {
	Create(ctx context.Context, data *model.Video) (*model.Video, error)
	Find(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.Video, error)
	List(
		ctx context.Context,
		filter *model.VideoFilter,
		paging *core.Paging,
		moreKeys ...string,
	) ([]model.Video, error)

	CreateVideoCategory(ctx context.Context, data *model.VideoCategory) (*model.VideoCategory, error)
	CreateVideoFile(ctx context.Context, data *model.VideoFile) (*model.VideoFile, error)
	GetVideoCategories(ctx context.Context, id uint64) ([]model.VideoCategory, error)
	GetVideoFile(ctx context.Context, id uint64) (*model.VideoFile, error)
}
