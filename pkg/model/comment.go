package model

type Comment struct {
	Base
	Content string `json:"content" gorm:"column:content;"`
	VideoID int64  `json:"video_id" gorm:"column:video_id;"`
	UserID  int64  `json:"user_id" gorm:"column:user_id;"`
}

func (Comment) TableName() string {
	return "comments"
}
