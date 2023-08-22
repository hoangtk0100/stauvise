package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xfrr/goffmpeg/transcoder"
)

func (server *Server) uploadVideo(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Bad Request")
		return
	}
	defer file.Close()

	filename := header.Filename
	ext := filepath.Ext(filename)
	baseFilename := strings.TrimSuffix(filename, ext)
	tempFilePath := filepath.Join("temp", filename)

	out, err := os.Create(tempFilePath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to copy file")
		return
	}

	hlsOutputDir := getHLSOutputDir(baseFilename)

	err = os.MkdirAll(hlsOutputDir, os.ModePerm)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to create output directory")
		return
	}

	trans := &transcoder.Transcoder{}
	err = trans.Initialize(tempFilePath, filepath.Join(hlsOutputDir, fmt.Sprintf("%s.m3u8", baseFilename)))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to initialize transcoder")
		return
	}

	done := trans.Run(false)
	err = <-done
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to convert video to HLS", err)
		return
	}

	ctx.String(http.StatusOK, "Upload and conversion successful")
}

func getHLSOutputDir(filename string) string {
	return filepath.Join("videos")
}
