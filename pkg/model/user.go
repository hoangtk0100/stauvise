package model

import "time"

type User struct {
	Base
	Username          string    `json:"username" gorm:"column:username;"`
	HashedPassword    string    `json:"-" gorm:"column:hashed_password;"`
	Salt              string    `json:"-" gorm:"column:salt;"`
	FullName          string    `json:"full_name" gorm:"column:full_name;"`
	Email             string    `json:"email" gorm:"column:email;"`
	PasswordChangedAt time.Time `json:"password_changed_at" gorm:"column:password_changed_at;"`
}

func (User) TableName() string {
	return "users"
}

type CreateUserParams struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type LoginParams struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}
