package model

import (
	"time"

	"github.com/hoangtk0100/stauvise/pkg/common"
)

type VideoFile struct {
	ID        int64           `json:"id" gorm:"column:id;" db:"id"`
	VideoID   int64           `json:"video_id" gorm:"column:video_id;"`
	Name      string          `json:"name" gorm:"column:name;"`
	Path      string          `json:"path" gorm:"column:path;"`
	Provider  common.Provider `json:"provider" gorm:"column:provider;default:LOCAL;"`
	CreatedAt *time.Time      `json:"created_at,omitempty" gorm:"column:created_at;" db:"created_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty" gorm:"column:deleted_at;" db:"deleted_at"`
}

func (VideoFile) TableName() string {
	return "video_files"
}
