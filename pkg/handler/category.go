package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param name body string true "Name" alphanum
// @Param description body string false "Description"
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /categories [post]
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

// GetCategories godoc
// @Summary Get all categories
// @Description Get all categories
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /categories [get]
func (server *Server) GetCategories(ctx *gin.Context) {
	categories, err := server.biz.Category().GetAll(ctx)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(categories))
}
