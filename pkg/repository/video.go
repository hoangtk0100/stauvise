package repository

import (
	"context"
	"github.com/hoangtk0100/stauvise/pkg/model"
	"strconv"
	"strings"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type videoRepo struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *videoRepo {
	return &videoRepo{db}
}

func (repo *videoRepo) Create(ctx context.Context, data *model.Video) (*model.Video, error) {
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

func (repo *videoRepo) Find(ctx context.Context, conds map[string]interface{}, moreInfo ...string) (*model.Video, error) {
	db := repo.db.Table(model.Video{}.TableName())
	for _, value := range moreInfo {
		db = db.Preload(value)
	}

	var video model.Video
	if err := db.Where(conds).First(&video).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &video, nil
}

func (repo *videoRepo) List(
	ctx context.Context,
	filter *model.VideoFilter,
	paging *core.Paging,
	moreKeys ...string,
) ([]model.Video, error) {
	var result []model.Video

	db := repo.db.
		Table(model.Video{}.TableName()).
		Where("status <> ?", common.VideoStatusDeleted)

	if f := filter; f != nil {
		if v := f.UserID; v != 0 {
			db = db.Where("owner_id = ?", v)
		}

		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		}
	}

	if err := db.Select("id").Count(&paging.Total).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	for _, value := range moreKeys {
		db = db.Preload(value)
	}

	if cursor := strings.TrimSpace(paging.FakeCursor); cursor != "" {
		id, err := strconv.ParseInt(cursor, 10, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		db = db.Where("id < ?", id)
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.
		Select("*").
		Order("id desc").
		Limit(paging.Limit).
		Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	size := len(result)
	if size > 0 {
		paging.NextCursor = strconv.FormatUint(result[size-1].ID, 10)
	}

	return result, nil
}

func (repo *videoRepo) CreateVideoCategory(ctx context.Context, data *model.VideoCategory) (*model.VideoCategory, error) {
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

func (repo *videoRepo) CreateVideoFile(ctx context.Context, data *model.VideoFile) (*model.VideoFile, error) {
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

func (repo *videoRepo) GetVideoCategories(ctx context.Context, id uint64) ([]model.VideoCategory, error) {
	db := repo.db.Table(model.VideoCategory{}.TableName()).Preload("Category")

	var result []model.VideoCategory
	if err := db.Where("video_id = ?", id).Find(&result).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (repo *videoRepo) GetVideoFile(ctx context.Context, id uint64) (*model.VideoFile, error) {
	var result model.VideoFile
	if err := repo.db.Where("video_id = ?", id).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &result, nil
}
