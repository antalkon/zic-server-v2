package handlers

import (
	"backend/internal/transport/rest/req"
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
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		code, msg := utils.InternalServerError("user ID is not valid UUID")
		return c.JSON(code, msg)
	}

	user, err := h.user.GetUserByID(userID)
	if err != nil {
		code, msg := utils.InternalServerError("failed to get user: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, user)
}
func (h *UserHandler) UpdateUser(c echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		code, msg := utils.InternalServerError("user ID is not valid UUID")
		return c.JSON(code, msg)
	}

	var userUpdate req.UpdateUserDataReq
	if err := c.Bind(&userUpdate); err != nil {
		code, msg := utils.NewError(http.StatusBadRequest, "failed to bind request: "+err.Error())
		return c.JSON(code, msg)
	}

	err := h.user.UpdateUser(userID, &userUpdate)
	if err != nil {
		code, msg := utils.InternalServerError("failed to update user: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) UpdatePassword(c echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		code, msg := utils.InternalServerError("user ID is not valid UUID")
		return c.JSON(code, msg)
	}

	var passwordUpdate req.UpdateDataPasswordReq
	if err := c.Bind(&passwordUpdate); err != nil {
		code, msg := utils.NewError(http.StatusBadRequest, "failed to bind request: "+err.Error())
		return c.JSON(code, msg)
	}

	err := h.user.UpdatePassword(userID, &passwordUpdate)
	if err != nil {
		code, msg := utils.InternalServerError("failed to update password: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.NoContent(http.StatusNoContent)
}
