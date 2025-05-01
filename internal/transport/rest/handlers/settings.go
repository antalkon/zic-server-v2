package handlers

import (
	"backend/internal/transport/rest/req"
	"backend/internal/transport/service"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type SettingsHandler struct {
	settings *service.SettingsService
	validate *validator.Validate
}

func NewSettingsHandler(settings *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		settings: settings,
		validate: validator.New(),
	}
}

// Echo handler
func (h *SettingsHandler) GetGeneralSettings(c echo.Context) error {
	settings, err := h.settings.GetGeneralSettings()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get settings",
		})
	}
	return c.JSON(http.StatusOK, settings)
}

func (h *SettingsHandler) UpdateGeneralSettings(c echo.Context) error {
	var req req.UpdateGeneralSettingsReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to bind request",
		})
	}

	// Пример: роль берётся из контекста, подставь свою логику
	userRole := c.Get("user_role").(string)
	if err := h.settings.UpdateGeneralSettings(&req, userRole); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Настройки успешно обновлены",
	})
}

func (h *SettingsHandler) GetTelegramSettings(c echo.Context) error {
	settings, err := h.settings.GetTelegramSettings()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get settings",
		})
	}
	return c.JSON(http.StatusOK, settings)
}

func (h *SettingsHandler) UpdateTelegramSettings(c echo.Context) error {
	var req req.UpdateTelegramSettingsReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to bind request",
		})
	}

	// Пример: роль берётся из контекста, подставь свою логику
	userRole := c.Get("user_role").(string)
	if err := h.settings.UpdateTelegramSettings(&req, userRole); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Настройки успешно обновлены",
	})
}
