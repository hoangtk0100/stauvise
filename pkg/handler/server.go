package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mdw "github.com/hoangtk0100/app-context/component/server/gin/middleware"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/docs"
	"github.com/hoangtk0100/stauvise/pkg/business"
	"github.com/hoangtk0100/stauvise/pkg/common"
	"github.com/hoangtk0100/stauvise/pkg/middleware"
	"github.com/hoangtk0100/stauvise/pkg/repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var (
	v1Prefix = "/api/v1"
)

type Server struct {
	config     *common.Config
	db         *gorm.DB
	tokenMaker core.TokenMakerComponent
	router     *gin.Engine
	ginsrv     core.GinComponent
	repo       repository.Repository
	biz        business.Business
}

func NewServer(config *common.Config) *Server {
	server := &Server{
		config:     config,
		db:         config.DB,
		tokenMaker: config.TokenMaker,
		ginsrv:     config.Gin,
		router:     config.Gin.GetRouter(),
	}

	server.repo = repository.NewRepository(server.db)
	server.biz = business.NewBusiness(server.repo, server.tokenMaker)

	server.setupRoutes()
	return server
}

func (server *Server) setupRoutes() {
	server.router.Use(mdw.Recovery(server.config.CTX), mdw.CORS(nil))
	v1 := server.router.Group(v1Prefix)
	authMdw := middleware.RequireAuth(server.repo.User(), server.GetTokenMaker())

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = v1Prefix

	auth := v1.Group("/auth")
	auth.POST("/login", server.Login)

	users := v1.Group("/users")
	users.POST("/register", server.Register)
	users.GET("/me", authMdw, server.GetProfile)

	categories := v1.Group("/categories").Use(authMdw)
	categories.POST("", server.CreateCategory)
	categories.GET("", server.GetCategories)

	videos := v1.Group("/videos")
	videos.StaticFS("/streams/", http.Dir("videos"))
	videos.POST("", authMdw, server.CreateVideo)
	videos.GET("", server.GetVideos)
	videos.GET("/:id", server.GetVideoDetails)
	videos.GET("/:id/segments", server.GetSegments)
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
