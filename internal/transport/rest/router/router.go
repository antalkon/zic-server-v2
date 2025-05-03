package router

import (
	"backend/internal/repository"
	"backend/internal/transport/middleware"
	"backend/internal/transport/rest/handlers"
	"backend/internal/transport/service"
	"backend/internal/tunel_service"
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
	userRepo := repository.NewUserRepository(ddbb)
	roomRepo := repository.NewRoomRepository(ddbb)
	computerRepo := repository.NewComputerRepository(ddbb)
	fristRepo := repository.NewFristRepository(ddbb)
	settingsRepo := repository.NewSettingsRepository(ddbb)
	tunelRpo := repository.NewTunelRepository(ddbb)
	actionsRepo := repository.NewActionsRepository(ddbb)

	authService := service.NewAuthService(authRepo)
	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo)
	roomService := service.NewRoomService(roomRepo)
	computerService := service.NewComputerService(computerRepo)
	fristService := service.NewFristService(fristRepo)
	settingsService := service.NewSettingsService(settingsRepo)
	tunelService := tunel_service.NewTunelService(tunelRpo, cache) // костыль
	actionsService := service.NewActionsService(actionsRepo, cache)

	authHandler := handlers.NewAuthHandler(authService)
	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)
	roomHandler := handlers.NewRoomHandler(roomService)
	computerHandler := handlers.NewComputerHandler(computerService)
	fristHandler := handlers.NewFristHandler(fristService)
	settingsHandler := handlers.NewSettingsHandler(settingsService)
	tunelHandler := handlers.NewTunelHandler(tunelService)
	actionsHandler := handlers.NewActionsHandler(actionsService)

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

	if config.ServiceGet().Server.Frist {
		frist := api.Group("/first")
		{
			frist.POST("/licenze/activate", fristHandler.ActivateLicenze)
			frist.POST("/form/org", fristHandler.OrgFormation)
			frist.POST("/form/api", fristHandler.ApiFormation)

		}
	}
	websocket := api.Group("/ws")
	{
		websocket.GET("/tunnel/:uuid", tunelHandler.HandleTunnel(cache)) // ✅ вот так правильно
	}

	data := api.Group("/data")
	data.Use(authMiddleware.AuthRequired(), middleware.LicenseRequired())
	{
		team := data.Group("/team")
		{
			team.POST("/role", teamHandler.CreateRole)
			team.GET("/role", teamHandler.GetAllRoles)     // Доработать + стат
			team.GET("/role/:id", teamHandler.GetRoleByID) // Доработать + стат
			team.PUT("/role/:id", teamHandler.UpdateRole)
			team.DELETE("/role/:id", teamHandler.DeleteRole) // Доработать + проверка на наличие юзеров
			team.POST("/user", teamHandler.CreateUser)
			team.GET("/user", teamHandler.GetAllUsers)     // Доработать + стат
			team.GET("/user/:id", teamHandler.GetUserByID) // Доработать + стат
			team.PUT("/user/:id", teamHandler.UpdateUser)
			team.DELETE("/user/:id", teamHandler.DeleteUser) // Доработать + проверка на наличие юзеров
			team.POST("/user/password/:id", teamHandler.UpdatePassword)
		}
		user := data.Group("/user")
		{
			user.GET("/data", userHandler.GetUserByID) // Получение данных пользователя, в т.ч. ролей и прав
			user.PUT("/data", userHandler.UpdateUser)  // Изменение данных пользователя
			user.POST("/password", userHandler.UpdatePassword)
		}
		room := data.Group("/room")
		{
			room.POST("", roomHandler.CreateRoom)
			room.GET("", roomHandler.GetAllRooms) // Доработать + стат
			room.GET("/:id", roomHandler.GetRoomByID)
			room.PUT("/:id", roomHandler.UpdateRoom)    // Доработать + стат
			room.DELETE("/:id", roomHandler.DeleteRoom) // Доработать + проверка на наличие юзеров
			room.GET("/computers/:id", roomHandler.GetRoomComputers)
			room.GET("/status/:id", roomHandler.GetRoomStatus) // Доработать + стат
		}
		computer := data.Group("/pc")
		{
			computer.POST("", computerHandler.CreateComputer)
			computer.GET("", computerHandler.GetAllComputers) // Доработать + стат
			computer.GET("/:id", computerHandler.GetComputerByID)
			computer.PUT("/:id", computerHandler.UpdateComputer)           // Доработать + стат
			computer.GET("/room/:id", computerHandler.GetRoomComputersAll) // Доработать + стат
			computer.DELETE("/:id", computerHandler.DeleteComputer)        // Доработать + проверка на наличие юзеров
		}
		settings := data.Group("/settings")
		{
			settings.GET("/general", settingsHandler.GetGeneralSettings) // Получение общих настроек
			settings.PUT("/general", settingsHandler.UpdateGeneralSettings)
			settings.GET("/telegram", settingsHandler.GetTelegramSettings) // Получение настроек телеги
			settings.PUT("/telegram", settingsHandler.UpdateTelegramSettings)
			settings.GET("/api", settingsHandler.GetApiSettings) // Получение настроек API
			settings.PUT("/api", settingsHandler.UpdateApiSettings)
			settings.GET("/license", settingsHandler.GetLicenseSettings) // Получение настроек лицензии
			settings.PUT("/license", settingsHandler.UpdateLicenseSettings)
		}
	}
	actions := api.Group("/actions")
	actions.Use(authMiddleware.AuthRequired(), middleware.LicenseRequired())
	{
		actions.POST("/reboot", actionsHandler.SendReboot)         // Перезагрузка
		actions.POST("/shutdown", actionsHandler.SendShutdown)     // Выключение
		actions.POST("/block", actionsHandler.SendBlock)           // Блокировка
		actions.POST("/unblock", actionsHandler.SendUnblock)       // Разблокировка
		actions.POST("/lockscreen", actionsHandler.SendLockScreen) // Блокировка экрана
		actions.POST("/sendurl", actionsHandler.SendUrl)           // Отправка URL
		actions.POST("/sendmessage", actionsHandler.SendMessage)   // Отправка сообщения

	}

}
