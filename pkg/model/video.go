package model

import (
	"github.com/hoangtk0100/stauvise/pkg/common"
)

type Video struct {
	Base
	Title       string             `json:"title" gorm:"column:title;"`
	Description string             `json:"description" gorm:"column:description;"`
	Status      common.VideoStatus `json:"status" gorm:"column:status;default:ACTIVE;"`
	OldStatus   common.VideoStatus `json:"-" gorm:"column:old_status;default:null"`
	OwnerID     uint64             `json:"owner_id" gorm:"column:owner_id;"`
}

func (Video) TableName() string {
	return "videos"
}

type CreateVideoParams struct {
	Title       string   `form:"title" binding:"required"`
	Description string   `form:"description,omitempty"`
	CategoryIDs []uint64 `form:"category_ids" binding:"required"`
}

type VideoFilter struct {
	UserID uint64 `json:"user_id" form:"user_id"`
	Status string `json:"status" form:"status"`
}
