package repository

import "gorm.io/gorm"

type repository struct {
	db          *gorm.DB
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db:          db,
		userRepo:    NewUserRepository(db),
		sessionRepo: NewSessionRepository(db),
	}
}

func (repo *repository) User() UserRepository {
	return repo.userRepo
}

func (repo *repository) Session() SessionRepository {
	return repo.sessionRepo
}
