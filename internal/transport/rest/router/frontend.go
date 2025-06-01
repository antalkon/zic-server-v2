package router

import (
	"github.com/labstack/echo/v4"
)

// SetupFrontendRoutes — подключаем фронтенд без шаблонов
func SetupFrontendRoutes(e *echo.Echo) {
	// Подключаем статику
	e.Static("/static", "web/static")
	e.Static("/public", "web/public")

	// Роут для страницы логина
	e.GET("/login", func(c echo.Context) error {
		return c.File("web/public/pages/login.html")
	})

	// При желании добавь другие страницы
	e.GET("/", func(c echo.Context) error {
		return c.File("web/public/pages/index.html")
	})
}
