package model

import (
	"time"

	"github.com/hoangtk0100/stauvise/pkg/common"
)

type VideoFile struct {
	ID         uint64          `json:"id" gorm:"column:id;" db:"id"`
	VideoID    uint64          `json:"video_id" gorm:"column:video_id;"`
	OriginName string          `json:"origin_name" gorm:"column:origin_name;"`
	Name       string          `json:"name" gorm:"column:name;"`
	Path       string          `json:"path" gorm:"column:path;"`
	Format     string          `json:"format" gorm:"column:format;"`
	Provider   common.Provider `json:"provider" gorm:"column:provider;default:LOCAL;"`
	MaxSegment int             `json:"max_segment" gorm:"column:max_segment;"`
	CreatedAt  *time.Time      `json:"created_at,omitempty" gorm:"column:created_at;" db:"created_at"`
	DeletedAt  *time.Time      `json:"deleted_at,omitempty" gorm:"column:deleted_at;" db:"deleted_at"`
}

func (VideoFile) TableName() string {
	return "video_files"
}
