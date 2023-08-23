package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

func (server *Server) Login(ctx *gin.Context) {
	var req model.LoginParams
	if err := ctx.ShouldBind(&req); err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	tokens, err := server.biz.Auth().Login(ctx, &req)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(tokens))
}
