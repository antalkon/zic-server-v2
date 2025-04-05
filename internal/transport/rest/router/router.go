package router

import (
	"backend/internal/repository"
	"backend/internal/transport/middleware"
	"backend/internal/transport/rest/handlers"
	"backend/internal/transport/service"
	"backend/pkg/cache"
	"backend/pkg/config"
	"backend/pkg/db"
	"backend/pkg/logger"
	"backend/pkg/storage"

	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo, cfg *config.Config, log *logger.Logger, db *db.Database, cache *cache.RedisClient, s3 *storage.MinIOClient) {
	ddbb := db.DB

	authRepo := repository.NewAuthRepository(ddbb)
	teamRepo := repository.NewTeamnRepository(ddbb)

	authService := service.NewAuthService(authRepo)
	teamService := service.NewTeamService(teamRepo)

	authHandler := handlers.NewAuthHandler(authService)
	teamHandler := handlers.NewTeamHandler(teamService)

	authMiddleware := middleware.NewAuthMiddleware(authRepo)

	api := e.Group("/api/v1")
	api.GET("/ping", handlers.Ping)

	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUpUser)
		auth.POST("/sign-in", authHandler.SignInUser)
		auth.POST("/refresh-token", authHandler.RefreshToken)
		auth.POST("/sign-out", authHandler.SignOutUser)
	}

	data := api.Group("/data")
	data.Use(authMiddleware.AuthRequired())
	{
		team := data.Group("/team")
		{
			team.POST("/role", teamHandler.CreateRole)
			team.GET("/role", teamHandler.GetAllRoles)     // Доработать + стат
			team.GET("/role/:id", teamHandler.GetRoleByID) // Доработать + стат
			team.PUT("/role/:id", teamHandler.UpdateRole)
			team.DELETE("/role/:id", teamHandler.DeleteRole) // Доработать + проверка на наличие юзеров
		}
	}
}
