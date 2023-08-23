package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

func (server *Server) CreateCategory(ctx *gin.Context) {
	var req model.CreateCategoryParams
	if err := ctx.ShouldBind(&req); err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	category, err := server.biz.Category().Create(ctx, &req)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(category))
}

func (server *Server) GetCategories(ctx *gin.Context) {
	categories, err := server.biz.Category().GetAll(ctx)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(categories))
}
