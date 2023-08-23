package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id" gorm:"column:id;"`
	OwnerID      uint64    `json:"owner_id" gorm:"column:owner_id;"`
	RefreshToken string    `json:"refresh_token" gorm:"column:refresh_token;"`
	UserAgent    string    `json:"user_agent" gorm:"column:user_agent;"`
	ClientIp     string    `json:"client_ip" gorm:"column:client_ip;"`
	IsBlocked    bool      `json:"is_blocked" gorm:"column:is_blocked;"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"column:expires_at;"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;"`
}

func (Session) TableName() string {
	return "sessions"
}
