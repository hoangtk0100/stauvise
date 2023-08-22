package repository

import (
	"context"

	"github.com/hoangtk0100/stauvise/pkg/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type sessionRepo struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *sessionRepo {
	return &sessionRepo{db}
}

func (u *sessionRepo) Create(ctx context.Context, data *model.Session) (*model.Session, error) {
	tx := u.db.Begin()
	if err := tx.Create(&data).Error; err != nil {
		tx.Rollback()
		return nil, errors.WithStack(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.WithStack(err)
	}

	return data, nil
}
