package middleware

import (
	glob_config "backend/config"
	"backend/pkg/licenze"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LicenseRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Проверка глобального флага
			if glob_config.Licenze {
				return next(c)
			}

			// 2. Проверка через licenze.CheckLicenze
			if licenze.CheckLicenze() {
				// Обновляем глобальный флаг (опционально)
				glob_config.Licenze = true
				return next(c)
			}

			// 3. Лицензия невалидна
			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Лицензия не активна или истек срок действия",
			})
		}
	}
}
