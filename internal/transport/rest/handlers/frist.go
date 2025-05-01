package handlers

import (
	"backend/internal/transport/rest/req"
	"backend/internal/transport/rest/res"
	"backend/internal/transport/service"
	"backend/internal/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type FristHandler struct {
	frist    *service.FristService
	validate *validator.Validate
}

func NewFristHandler(frist *service.FristService) *FristHandler {
	return &FristHandler{
		frist:    frist,
		validate: validator.New(),
	}
}

func (h *FristHandler) ActivateLicenze(c echo.Context) error {
	// Получаем токен из тела запроса
	var req req.ActivateLicenzeReq
	if err := c.Bind(&req); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	if err := h.validate.Struct(req); err != nil {
		code, msg := utils.ValidationError()
		return c.JSON(code, msg)
	}

	// Проверяем лицензию
	if err := h.frist.ActivateLicenze(req.Token); err != nil {
		code, msg := utils.InternalServerError("failed to activate license: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.MessageRes{
		Message: "License activated successfully",
	})
}

func (h *FristHandler) OrgFormation(c echo.Context) error {
	// Получаем токен из тела запроса
	var req req.FristStartForm
	if err := c.Bind(&req); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	if err := h.validate.Struct(req); err != nil {
		code, msg := utils.ValidationError()
		return c.JSON(code, msg)
	}

	if err := h.frist.OrgFormation(&req); err != nil {
		code, msg := utils.InternalServerError("failed to create organization: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.MessageRes{
		Message: "License activated successfully",
	})
}

func (h *FristHandler) ApiFormation(c echo.Context) error {
	// Получаем токен из тела запроса
	var req req.FristStartAPI
	if err := c.Bind(&req); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	if err := h.validate.Struct(req); err != nil {
		code, msg := utils.ValidationError()
		return c.JSON(code, msg)
	}

	if err := h.frist.ApiFormation(&req); err != nil {
		code, msg := utils.InternalServerError("failed to create organization: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.MessageRes{
		Message: "License activated successfully",
	})
}
