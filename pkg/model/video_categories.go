package model

import "time"

type VideoCategory struct {
	ID         uint64     `json:"id" gorm:"column:id;" db:"id"`
	VideoID    uint64     `json:"video_id" gorm:"column:video_id;"`
	CategoryID uint64     `json:"category_id" gorm:"column:category_id;"`
	Category   Category   `json:"category" gorm:"foreignKey:CategoryID;"`
	CreatedAt  *time.Time `json:"created_at,omitempty" gorm:"column:created_at;" db:"created_at"`
}

func (VideoCategory) TableName() string {
	return "video_categories"
}
