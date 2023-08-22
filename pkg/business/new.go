package business

import (
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/stauvise/pkg/repository"
)

type business struct {
	repo    repository.Repository
	authBiz AuthBusiness
	userBiz UserBusiness
}

func NewBusiness(repo repository.Repository, tokenMaker token.TokenMaker) *business {
	return &business{
		repo:    repo,
		authBiz: NewAuthBusiness(repo, tokenMaker),
		userBiz: NewUserBusiness(repo),
	}
}

func (b *business) Auth() AuthBusiness {
	return b.authBiz
}

func (b *business) User() UserBusiness {
	return b.userBiz
}
