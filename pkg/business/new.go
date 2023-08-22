package business

import (
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/stauvise/pkg/repository"
)

type business struct {
	repo        repository.Repository
	authBiz     AuthBusiness
	userBiz     UserBusiness
	categoryBiz CategoryBusiness
}

func NewBusiness(repo repository.Repository, tokenMaker token.TokenMaker) *business {
	return &business{
		repo:        repo,
		authBiz:     NewAuthBusiness(repo, tokenMaker),
		userBiz:     NewUserBusiness(repo),
		categoryBiz: NewCategoryBusiness(repo),
	}
}

func (b *business) Auth() AuthBusiness {
	return b.authBiz
}

func (b *business) User() UserBusiness {
	return b.userBiz
}

func (b *business) Category() CategoryBusiness {
	return b.categoryBiz
}
