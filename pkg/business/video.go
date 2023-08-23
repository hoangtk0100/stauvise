package business

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/common"
	"github.com/hoangtk0100/stauvise/pkg/model"
	"github.com/hoangtk0100/stauvise/pkg/repository"
	"github.com/hoangtk0100/stauvise/pkg/util"
	"github.com/hoangtk0100/stauvise/pkg/validator"
)

var (
	ErrVideoTitleExisted     = errors.New("video title already existed")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCannotCreateVideo     = errors.New("cannot create video")
	ErrCannotCreateVideoFile = errors.New("cannot create video file")
	ErrCannotCreateFile      = errors.New("cannot create file")
	ErrCannotCopyFile        = errors.New("cannot copy file")
	ErrCannotCreateOutputDir = errors.New("cannot create output directory")
	ErrCannotInitTranscoder  = errors.New("cannot init transcoder")
	ErrCannotConvertMedia    = errors.New("cannot convert media")
	ErrCannotListVideos      = errors.New("cannot list videos")
	ErrVideoNotFound         = errors.New("video not found")
	ErrVideoSegmentsNotFound = errors.New("video segments not found")
	ErrCannotReadFile        = errors.New("cannot read file")
	ErrCannotRemoveTempFile  = errors.New("cannot remove temp file")
)

type videoBusiness struct {
	repo repository.Repository
}

func NewVideoBusiness(repo repository.Repository) *videoBusiness {
	return &videoBusiness{repo: repo}
}

type categoryResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type videoResponse struct {
	ID             uint64             `json:"id"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	FileName       string             `json:"file_name"`
	OriginFileName string             `json:"origin_file_name"`
	Path           string             `json:"path"`
	MaxSegment     int                `json:"max_segment"`
	Categories     []categoryResponse `json:"categories"`
	CreatedAt      *time.Time         `json:"created_at,omitempty"`
}

func newVideoResponse(video *model.Video, categories []model.VideoCategory, videoFile *model.VideoFile) videoResponse {
	var categoriesRsp []categoryResponse
	for _, itemIdx := range categories {
		categoriesRsp = append(categoriesRsp, categoryResponse{
			ID:   itemIdx.Category.ID,
			Name: itemIdx.Category.Name,
		})
	}

	return videoResponse{
		ID:             video.ID,
		Title:          video.Title,
		Description:    video.Description,
		FileName:       videoFile.Name,
		OriginFileName: videoFile.OriginName,
		Path:           videoFile.Path,
		MaxSegment:     videoFile.MaxSegment,
		Categories:     categoriesRsp,
		CreatedAt:      video.CreatedAt,
	}
}

func (biz *videoBusiness) Create(ctx context.Context, data *model.CreateVideoParams, fileHeader *multipart.FileHeader) (uint64, error) {
	requesterID, err := common.GetRequesterID(ctx)
	if err != nil {
		return 0, core.ErrUnauthorized.WithDebug(err.Error())
	}

	existedVideo, _ := biz.repo.Video().Find(ctx, map[string]interface{}{"title": data.Title})
	if existedVideo != nil {
		return 0, core.ErrConflict.WithError(ErrVideoTitleExisted.Error())
	}

	categories, err := biz.repo.Category().GetByIDs(ctx, data.CategoryIDs)
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return 0, core.ErrNotFound.WithError(ErrCategoryNotFound.Error())
		}

		return 0, core.ErrInternalServerError.
			WithError(ErrCategoryNotFound.Error()).
			WithDebug(err.Error())
	}

	videoFileParams, err := biz.processVideoFile(ctx, fileHeader)
	if err != nil {
		return 0, err
	}

	video, err := biz.repo.Video().Create(ctx, &model.Video{
		Title:       data.Title,
		Description: data.Description,
		OwnerID:     requesterID,
	})

	if err != nil {
		return 0, core.ErrInternalServerError.
			WithError(ErrCannotCreateVideo.Error()).
			WithDebug(err.Error())
	}

	for idx := range categories {
		_, err = biz.repo.Video().CreateVideoCategory(ctx, &model.VideoCategory{
			VideoID:    video.ID,
			CategoryID: categories[idx].ID,
		})

		if err != nil {
			return 0, core.ErrInternalServerError.
				WithError(ErrCannotCreateVideo.Error()).
				WithDebug(err.Error())
		}
	}

	videoFileParams.VideoID = video.ID
	_, err = biz.repo.Video().CreateVideoFile(ctx, videoFileParams)
	if err != nil {
		return 0, core.ErrInternalServerError.
			WithError(ErrCannotCreateVideoFile.Error()).
			WithDebug(err.Error())
	}

	return video.ID, nil
}

func (biz *videoBusiness) processVideoFile(ctx context.Context, fileHeader *multipart.FileHeader) (data *model.VideoFile, err error) {
	fileName := fileHeader.Filename
	ext := filepath.Ext(fileName)
	err = validator.ValidateFileExt(ext)
	if err != nil {
		err = core.ErrBadRequest.
			WithError(err.Error())
		return
	}

	folder, newFileName := getNewFileName(ext)
	dir := getDst(folder, folder)
	dst := fmt.Sprintf("%s.m3u8", dir)
	tempFilePath := filepath.Join("temp", newFileName)

	out, err := os.Create(tempFilePath)
	if err != nil {
		err = core.ErrInternalServerError.
			WithError(ErrCannotCreateFile.Error()).
			WithDebug(err.Error())
		return
	}
	defer out.Close()

	file, err := fileHeader.Open()
	if err != nil {
		err = core.ErrBadRequest.
			WithError(ErrCannotReadFile.Error()).
			WithDebug(err.Error())
		return
	}
	defer file.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		err = core.ErrInternalServerError.
			WithError(ErrCannotCopyFile.Error()).
			WithDebug(err.Error())
		return
	}

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		err = core.ErrInternalServerError.
			WithError(ErrCannotCreateOutputDir.Error()).
			WithDebug(err.Error())
		return
	}

	if err = util.ToHLS(ctx, tempFilePath, dst, dir); err != nil {
		err = core.ErrInternalServerError.
			WithError(ErrCannotConvertMedia.Error()).
			WithDebug(err.Error())
		return
	}

	_ = os.Remove(dst)
	err = os.Rename(fmt.Sprintf("%s.tmp", dst), dst)
	if err != nil {
		err = core.ErrInternalServerError.
			WithDebug(ErrCannotRemoveTempFile.Error())
		return
	}

	maxSegment, err := util.GetMaxSegmentNumber(dst)
	if err != nil {
		err = core.ErrInternalServerError.
			WithError(ErrCannotConvertMedia.Error()).
			WithDebug(err.Error())
		return
	}

	err = os.Remove(tempFilePath)
	if err != nil {
		err = core.ErrInternalServerError.
			WithDebug(ErrCannotRemoveTempFile.Error())
		return
	}

	params := &model.VideoFile{
		OriginName: fileName,
		Name:       newFileName,
		Path:       fmt.Sprintf("%s/%s.m3u8", folder, folder),
		Format:     ext,
		MaxSegment: maxSegment,
		Provider:   common.ProviderLocal,
	}

	return params, nil
}

func (biz *videoBusiness) List(ctx context.Context, filter *model.VideoFilter, paging *core.Paging) ([]model.Video, error) {
	videos, err := biz.repo.Video().List(ctx, filter, paging)
	if err != nil {
		return nil, core.ErrInternalServerError.
			WithError(ErrCannotListVideos.Error()).
			WithDebug(err.Error())
	}

	return videos, nil
}

func (biz *videoBusiness) GetByID(ctx context.Context, id uint64) (interface{}, error) {
	video, err := biz.repo.Video().Find(ctx, map[string]interface{}{"id": id})
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return nil, core.ErrNotFound.WithError(ErrVideoNotFound.Error())
		}

		return nil, core.ErrInternalServerError.
			WithError(ErrVideoNotFound.Error()).
			WithDebug(err.Error())
	}

	videoCategories, err := biz.repo.Video().GetVideoCategories(ctx, id)
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return nil, core.ErrNotFound.WithError(ErrVideoNotFound.Error())
		}

		return nil, core.ErrInternalServerError.
			WithError(ErrVideoNotFound.Error()).
			WithDebug(err.Error())
	}

	videoFile, err := biz.repo.Video().GetVideoFile(ctx, id)
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return nil, core.ErrNotFound.WithError(ErrVideoNotFound.Error())
		}

		return nil, core.ErrInternalServerError.
			WithError(ErrVideoNotFound.Error()).
			WithDebug(err.Error())
	}

	rsp := newVideoResponse(video, videoCategories, videoFile)
	return rsp, nil
}

func paginateFromRange(max int, paging *core.Paging) (from, to int) {
	from = (paging.Page - 1) * paging.Limit
	to = paging.Page * paging.Limit
	if to > max {
		to = max + 1
	}

	return
}

func (biz *videoBusiness) GetSegments(ctx context.Context, id uint64, paging *core.Paging) (interface{}, error) {
	data, err := biz.repo.Video().GetVideoFile(ctx, id)
	if err != nil {
		if core.ErrNotFound.Is(err) {
			return nil, core.ErrNotFound.WithError(ErrVideoSegmentsNotFound.Error())
		}

		return nil, core.ErrInternalServerError.
			WithError(ErrVideoSegmentsNotFound.Error()).
			WithDebug(err.Error())
	}

	paging.Total = int64(data.MaxSegment + 1)
	from, to := paginateFromRange(data.MaxSegment, paging)
	rsp := map[int]string{}

	for idx := from; idx < to; idx++ {
		rsp[idx] = getSegmentDst(data.Path, ".m3u8", ".ts", idx)
	}

	return rsp, nil
}

func getNewFileName(ext string) (string, string) {
	uid := time.Now().UTC().UnixNano()
	return fmt.Sprintf("%d", uid), fmt.Sprintf("%d%s", uid, ext)
}

func getDst(folder, fileName string) string {
	return fmt.Sprintf("videos/%s/%s", folder, fileName)
}

func getSegmentDst(dst, oldExt, ext string, idx int) string {
	suffix := fmt.Sprintf("_%d%s", idx, ext)

	return strings.ReplaceAll(dst, oldExt, suffix)
}
