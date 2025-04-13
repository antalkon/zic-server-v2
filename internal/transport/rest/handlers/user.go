package handlers

import (
	"backend/internal/transport/service"
	"backend/internal/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	user     *service.UserService
	validate *validator.Validate
}

func NewUserHandler(user *service.UserService) *UserHandler {
	return &UserHandler{
		user:     user,
		validate: validator.New(),
	}
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	userIdStr, ok := c.Get("user_id").(string)
	if !ok || userIdStr == "" {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	user, err := h.user.GetUserByID(userID)
	if err != nil {
		code, msg := utils.InternalServerError("failed to get user: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, user)
}
