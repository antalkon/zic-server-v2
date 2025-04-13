package handlers

import (
	"backend/internal/models"
	"backend/internal/transport/rest/req"
	"backend/internal/transport/service"
	"backend/internal/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ComputerHandler struct {
	computer *service.ComputerService
	validate *validator.Validate
}

func NewComputerHandler(computer *service.ComputerService) *ComputerHandler {
	return &ComputerHandler{
		computer: computer,
		validate: validator.New(),
	}
}

func (h *ComputerHandler) CreateComputer(c echo.Context) error {
	var computer models.Computer

	// Привязка JSON к модели
	if err := c.Bind(&computer); err != nil {
		code, msg := utils.NewError(http.StatusBadRequest, "failed to bind request: "+err.Error())
		return c.JSON(code, msg)
	}

	// Валидация модели
	if err := h.validate.Struct(computer); err != nil {
		code, msg := utils.NewError(http.StatusBadRequest, "validation failed: "+err.Error())
		return c.JSON(code, msg)
	}

	// Получение роли пользователя из контекста
	userRole, ok := c.Get("user_role").(string)
	if !ok || userRole == "" {
		code, msg := utils.UnauthorizedError()
		return c.JSON(code, msg)
	}

	// Вызов бизнес-логики для создания компьютера
	jwt, err := h.computer.CreateComputer(&computer, userRole)
	if err != nil {
		code, msg := utils.NewError(http.StatusInternalServerError, "failed to create computer: "+err.Error())
		return c.JSON(code, msg)
	}

	// Возврат успешного ответа с JWT
	return c.JSON(http.StatusCreated, map[string]string{
		"token": jwt,
	})
}

func (h *ComputerHandler) GetAllComputers(c echo.Context) error {
	computers, err := h.computer.GetAllComputers()
	if err != nil {
		code, msg := utils.NewError(http.StatusInternalServerError, "failed to get computers: "+err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, computers)
}

func (h *ComputerHandler) GetComputerByID(c echo.Context) error {
	computerID := c.Param("id")

	computer, err := h.computer.GetComputerByID(computerID)
	if err != nil {
		code, msg := utils.NewError(http.StatusInternalServerError, "failed to get computer: "+err.Error())
		return c.JSON(code, msg)
	}

	if computer == nil {
		code, msg := utils.NewError(http.StatusNotFound, "computer not found")
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, computer)
}

func (h *ComputerHandler) GetRoomComputersAll(c echo.Context) error {
	roomIDStr := c.Param("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	computers, err := h.computer.GetRoomComputersAll(roomID)
	if err != nil {
		code, msg := utils.InternalServerError("failed to get computers: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, computers)
}

func (h *ComputerHandler) UpdateComputer(c echo.Context) error {
	idStr := c.Param("id")
	computerID, err := uuid.Parse(idStr)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}
	userRole, ok := c.Get("user_role").(string)
	if !ok || userRole == "" {
		code, msg := utils.UnauthorizedError()
		return c.JSON(code, msg)
	}

	var update req.UpdateComputerReq
	if err := c.Bind(&update); err != nil {
		code, msg := utils.NewError(http.StatusBadRequest, "failed to bind request: "+err.Error())
		return c.JSON(code, msg)
	}

	err = h.computer.UpdateComputer(computerID, &update, userRole)
	if err != nil {
		code, msg := utils.InternalServerError("failed to update computer: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ComputerHandler) DeleteComputer(c echo.Context) error {
	computerIDStr := c.Param("id")
	computerID, err := uuid.Parse(computerIDStr)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	userRole, ok := c.Get("user_role").(string)
	if !ok || userRole == "" {
		code, msg := utils.UnauthorizedError()
		return c.JSON(code, msg)
	}

	err = h.computer.DeleteComputer(computerID, userRole)
	if err != nil {
		code, msg := utils.InternalServerError("failed to delete computer: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.NoContent(http.StatusNoContent)
}
