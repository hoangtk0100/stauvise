package repository

import (
	"context"

	"github.com/hoangtk0100/stauvise/pkg/model"
)

type Repository interface {
	User() UserRepository
	Session() SessionRepository
	Category() CategoryRepository
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
	GetAll(ctx context.Context) ([]model.Category, error)
}
