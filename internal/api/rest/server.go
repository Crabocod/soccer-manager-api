package rest

import (
	"soccer_manager_service/internal/api/rest/handlers"
	"soccer_manager_service/internal/api/rest/middleware"
	"soccer_manager_service/internal/usecase"
	i18nPkg "soccer_manager_service/pkg/i18n"
	"soccer_manager_service/pkg/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "soccer_manager_service/internal/api/rest/swagger/docs"
)

// @title Soccer Manager API
// @version 1.0
// @description REST API for Soccer Manager Service

// @contact.name API Support
// @contact.email support@soccermanager.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

type Server struct {
	router      *gin.Engine
	jwtManager  *jwt.Manager
	usecase     *usecase.Service
	logger      *zap.Logger
	i18nManager *i18nPkg.Manager
}

func NewServer(jwtManager *jwt.Manager, usecase *usecase.Service, logger *zap.Logger, i18nManager *i18nPkg.Manager) *Server {
	router := gin.Default()

	s := &Server{
		router:      router,
		jwtManager:  jwtManager,
		usecase:     usecase,
		logger:      logger,
		i18nManager: i18nManager,
	}

	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	authHandler := handlers.NewAuthHandler(s.usecase.Auth, s.logger)
	teamHandler := handlers.NewTeamHandler(s.usecase.Team, s.logger)
	playerHandler := handlers.NewPlayerHandler(s.usecase.Player, s.logger)
	transferHandler := handlers.NewTransferHandler(s.usecase.Transfer, s.logger)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := s.router.Group("/api/v1")
	api.Use(middleware.I18nMiddleware(s.i18nManager))
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		authMiddleware := middleware.Auth(s.jwtManager)

		team := api.Group("/team")
		team.Use(authMiddleware)
		{
			team.GET("", teamHandler.GetMyTeam)
			team.PATCH("", teamHandler.UpdateTeam)
		}

		players := api.Group("/players")
		players.Use(authMiddleware)
		{
			players.PATCH("/:id", playerHandler.UpdatePlayer)
			players.POST("/:id/transfer", transferHandler.ListPlayer)
		}

		transfers := api.Group("/transfers")
		transfers.Use(authMiddleware)
		{
			transfers.GET("", transferHandler.GetTransferList)
			transfers.POST("/:id/buy", transferHandler.BuyPlayer)
		}
	}
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
