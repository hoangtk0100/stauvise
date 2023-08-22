package main

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

func main() {
	router := gin.Default()
	router.Use(CORS(nil))

	router.StaticFS("/sources/", http.Dir("videos/"))
	router.POST("/upload", uploadVideo)

	router.Run(":8080")
}

var (
	defaultHeaders = map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, PATH, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":     "Origin, Authorization, Content-Type",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Max-Age":           "86400",
	}
)

func CORS(headers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Override default headers
		if len(headers) != 0 {
			for key, value := range headers {
				defaultHeaders[key] = value
			}
		}

		// Set headers
		for key, value := range defaultHeaders {
			c.Writer.Header().Set(key, value)
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func uploadVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}
	defer file.Close()

	filename := header.Filename
	ext := filepath.Ext(filename)
	baseFilename := strings.TrimSuffix(filename, ext)
	tempFilePath := filepath.Join("temp", filename)

	out, err := os.Create(tempFilePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to copy file")
		return
	}

	hlsOutputDir := getHLSOutputDir(baseFilename)

	err = os.MkdirAll(hlsOutputDir, os.ModePerm)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create output directory")
		return
	}

	trans := &transcoder.Transcoder{}
	err = trans.Initialize(tempFilePath, filepath.Join(hlsOutputDir, fmt.Sprintf("%s.m3u8", baseFilename)))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to initialize transcoder")
		return
	}

	done := trans.Run(false)
	err = <-done
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to convert video to HLS", err)
		return
	}

	c.String(http.StatusOK, "Upload and conversion successful")
}

func getHLSOutputDir(filename string) string {
	return filepath.Join("videos")
}
