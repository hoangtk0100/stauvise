package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

func (server *Server) CreateVideo(ctx *gin.Context) {
	var req model.CreateVideoParams
	if err := ctx.ShouldBind(&req); err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	id, err := server.biz.Video().Create(ctx, &req, fileHeader)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(id))
}

func (server *Server) GetVideos(ctx *gin.Context) {
	var queryString struct {
		core.Paging
		model.VideoFilter
	}

	if err := ctx.ShouldBind(&queryString); err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	queryString.Paging.Process()

	videos, err := server.biz.Video().List(ctx, &queryString.VideoFilter, &queryString.Paging)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewResponse(videos, queryString.Paging, queryString.VideoFilter))
}

type videoIDRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

func (server *Server) GetVideoDetails(ctx *gin.Context) {
	var reqID videoIDRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	video, err := server.biz.Video().GetByID(ctx, reqID.ID)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewDataResponse(video))
}

func (server *Server) GetSegments(ctx *gin.Context) {
	var reqID videoIDRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	var paging core.Paging
	if err := ctx.ShouldBind(&paging); err != nil {
		core.ErrorResponse(ctx, core.ErrBadRequest.WithError(err.Error()))
		return
	}

	segments, err := server.biz.Video().GetSegments(ctx, reqID.ID, &paging)
	if err != nil {
		core.ErrorResponse(ctx, err)
		return
	}

	core.SuccessResponse(ctx, core.NewResponse(segments, paging, map[string]uint64{"video_id": reqID.ID}))
}
