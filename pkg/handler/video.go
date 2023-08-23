package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/model"
)

// CreateVideo godoc
// @Summary Create a new video
// @Description Create a new video
// @Tags video
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param title formData string true "Title"
// @Param description formData string false "Description"
// @Param category_ids formData []uint64 true "CategoryIDs"
// @Param file formData file true "Media file"
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /videos [post]
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

// GetVideos godoc
// @Summary List videos by filters
// @Description List videos by filters
// @Tags video
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param user_id query int false "UserID"
// @Param status query string false "Status"
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /videos [get]
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

// GetVideoDetails godoc
// @Summary Get video details
// @Description Get video details
// @Tags video
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /videos/{id} [get]
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

// GetSegments godoc
// @Summary Get video segments
// @Description Get video segments
// @Tags video
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} successResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /videos/{id}/segments [get]
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
