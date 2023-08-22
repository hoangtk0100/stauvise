package repository

import "gorm.io/gorm"

type repository struct {
	db          *gorm.DB
	userRepo    UserRepository
	sessionRepo SessionRepository
	cateRepo    CategoryRepository
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db:          db,
		userRepo:    NewUserRepository(db),
		sessionRepo: NewSessionRepository(db),
		cateRepo:    NewCategoryRepository(db),
	}
}

func (repo *repository) User() UserRepository {
	return repo.userRepo
}

func (repo *repository) Session() SessionRepository {
	return repo.sessionRepo
}

func (repo *repository) Category() CategoryRepository {
	return repo.cateRepo
}
