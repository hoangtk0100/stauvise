package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

// Register godoc
// @Summary Register
// @Description Register
// @Tags user
// @Accept json
// @Produce json
// @Param username body string true "Username" alphanum
// @Param password body string true "Password" minlength(8)
// @Param full_name body string true "Full Name" minlength(8)
// @Param email body string true "Email" minlength(8)
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /users/register [post]
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

// GetProfile godoc
// @Summary Get current user profile
// @Description Get current user profile
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /users/me [get]
func (server *Server) GetProfile(ctx *gin.Context) {
	user, err := server.biz.User().GetProfile(ctx)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(user))
}
