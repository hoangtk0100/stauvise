package repository

import "gorm.io/gorm"

type repository struct {
	db           *gorm.DB
	userRepo     UserRepository
	sessionRepo  SessionRepository
	categoryRepo CategoryRepository
	videoRepo    VideoRepository
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db:           db,
		userRepo:     NewUserRepository(db),
		sessionRepo:  NewSessionRepository(db),
		categoryRepo: NewCategoryRepository(db),
		videoRepo:    NewVideoRepository(db),
	}
}

func (repo *repository) User() UserRepository {
	return repo.userRepo
}

func (repo *repository) Session() SessionRepository {
	return repo.sessionRepo
}

func (repo *repository) Category() CategoryRepository {
	return repo.categoryRepo
}

func (repo *repository) Video() VideoRepository {
	return repo.videoRepo
}
