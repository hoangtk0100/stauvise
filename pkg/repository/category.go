package repository

import (
	"context"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepo {
	return &categoryRepo{db}
}

func (repo *categoryRepo) GetByIDs(ctx context.Context, ids []uint64) ([]model.Category, error) {
	var result []model.Category
	if err := repo.db.
		Table(model.Category{}.TableName()).
		Where("id in (?)", ids).
		Find(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (repo *categoryRepo) GetByName(ctx context.Context, name string) (*model.Category, error) {
	var result model.Category
	if err := repo.db.Where("name = ?", name).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func (repo *categoryRepo) Create(ctx context.Context, data *model.Category) (*model.Category, error) {
	tx := repo.db.Begin()
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

func (repo *categoryRepo) GetAll(ctx context.Context) ([]model.Category, error) {
	var result []model.Category
	if err := repo.db.Table(model.Category{}.TableName()).
		Order("id desc").
		Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
