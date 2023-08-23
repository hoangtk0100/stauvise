package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

func (server *Server) Register(ctx *gin.Context) {
	var req model.CreateUserParams
	if err := ctx.ShouldBind(&req); err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	if _, err := server.biz.User().Register(ctx, &req); err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(true))
}

func (server *Server) GetProfile(ctx *gin.Context) {
	user, err := server.biz.User().GetProfile(ctx)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(user))
}
