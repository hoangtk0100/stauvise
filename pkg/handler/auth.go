package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

type successResponse struct {
	Data   interface{} `json:"data"`
	Paging interface{} `json:"paging,omitempty"`
	Extra  interface{} `json:"extra,omitempty"`
}

type errorResponse struct {
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Error  string `json:"message"`
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param username body string true "Username" alphanum
// @Param password body string true "Password" minlength(8)
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /auth/login [post]
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
