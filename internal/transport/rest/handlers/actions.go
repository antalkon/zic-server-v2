package handlers

import (
	"backend/internal/transport/rest/req"
	"backend/internal/transport/service"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ActionsHandler struct {
	actions  *service.ActionsService
	validate *validator.Validate
}

func NewActionsHandler(actions *service.ActionsService) *ActionsHandler {
	return &ActionsHandler{
		actions:  actions,
		validate: validator.New(),
	}
}

// Стандартные обработчик добавить!!!!!!!!!!!!
func (h *ActionsHandler) SendReboot(c echo.Context) error {
	var r req.SendRebootReq
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}
	if err := h.validate.Struct(r); err != nil {
		return echo.NewHTTPError(400, "Validation error", err.Error())
	}
	if err := h.actions.SendReboot(r.ComputerID, r.Delay); err != nil {
		return echo.NewHTTPError(500, "Failed to send reboot command", err.Error())
	}
	return c.JSON(200, map[string]string{"status": "success"})
}
