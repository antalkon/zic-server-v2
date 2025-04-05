package handlers

import (
	"backend/internal/models"
	"backend/internal/transport/rest/res"
	"backend/internal/transport/service"
	"backend/internal/utils"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type TeamHandler struct {
	team     *service.TeamService
	validate *validator.Validate
}

func NewTeamHandler(team *service.TeamService) *TeamHandler {
	return &TeamHandler{
		team:     team,
		validate: validator.New(),
	}
}

func (h *TeamHandler) CreateRole(c echo.Context) error {
	var role models.Role

	if err := c.Bind(&role); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	if err := h.validate.Struct(role); err != nil {
		code, msg := utils.ValidationError()
		return c.JSON(code, msg)
	}

	err := h.team.CreateRole(&role)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			code, msg := utils.ConflictCustomError("Role with this name already exists")
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to create role: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusCreated, res.CreateRoleRes{
		Message: "Role created successfully",
	})
}
