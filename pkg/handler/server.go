package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mdw "github.com/hoangtk0100/app-context/component/server/gin/middleware"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/app-context/core"
	"github.com/hoangtk0100/stauvise/pkg/business"
	"github.com/hoangtk0100/stauvise/pkg/common"
	"github.com/hoangtk0100/stauvise/pkg/middleware"
	"github.com/hoangtk0100/stauvise/pkg/repository"
	"gorm.io/gorm"
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
	router := server.router
	router.Use(mdw.Recovery(server.config.CTX), mdw.CORS(nil))
	authMdw := middleware.RequireAuth(server.repo.User(), server.GetTokenMaker())

	v1 := router.Group("/v1")

	auth := v1.Group("/auth")
	auth.POST("/login", server.login)

	users := v1.Group("/users")
	users.POST("/register", server.register)
	users.Use(authMdw)
	users.GET("/me", server.getProfile)

	categories := v1.Group("/categories").Use(authMdw)
	categories.POST("/", server.createCategory)
	categories.GET("/", server.getCategories)

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
