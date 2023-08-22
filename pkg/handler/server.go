package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/component/server/gin/middleware"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/common"
)

type Server struct {
	config     *common.Config
	tokenMaker core.TokenMakerComponent
	router     *gin.Engine
	ginsrv     core.GinComponent
}

func NewServer(config *common.Config) *Server {
	server := &Server{
		config:     config,
		tokenMaker: config.TokenMaker,
		ginsrv:     config.Gin,
		router:     config.Gin.GetRouter(),
	}

	server.setupRoutes()
	return server
}

func (server *Server) setupRoutes() {
	router := server.router
	router.Use(middleware.CORS(nil))

	v1 := router.Group("/v1")

	v1.StaticFS("/sources/", http.Dir("videos/"))
	v1.POST("/upload", server.uploadVideo)

	server.router = router
}

func (server *Server) GetRouter() *gin.Engine {
	return server.router
}

func (server *Server) GetTokenMaker() token.TokenMaker {
	return server.tokenMaker
}

func (server *Server) RunDBMigration() {
	server.config.DBMigrator.MigrateUp()
}

func (server *Server) Start() {
	server.ginsrv.StartGracefully()
}
