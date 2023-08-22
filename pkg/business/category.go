package business

import (
	"context"
	"errors"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
	"github.com/hoangtk0100/stauvise/pkg/repository"
)

var (
	ErrCategoryExisted      = errors.New("category existed")
	ErrCannotCreateCategory = errors.New("cannot create category")
	ErrCannotGetCategories  = errors.New("cannot get categories")
)

type categoryBusiness struct {
	repo repository.Repository
}

func NewCategoryBusiness(repo repository.Repository) *categoryBusiness {
	return &categoryBusiness{repo: repo}
}

func (biz *categoryBusiness) Create(ctx context.Context, data *model.CreateCategoryParams) (*model.Category, error) {
	existedCategory, err := biz.repo.Category().GetByName(ctx, data.Name)
	if err == nil && existedCategory != nil {
		return nil, core.ErrConflict.WithError(ErrCategoryExisted.Error())
	}

	params := &model.Category{
		Name:        data.Name,
		Description: data.Description,
	}

	category, err := biz.repo.Category().Create(ctx, params)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrCannotCreateCategory.Error()).
			WithDebug(err.Error())
	}

	return category, nil
}

func (biz *categoryBusiness) GetAll(ctx context.Context) ([]model.Category, error) {
	categories, err := biz.repo.Category().GetAll(ctx)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrCannotGetCategories.Error()).
			WithDebug(err.Error())
	}

	return categories, nil
}
