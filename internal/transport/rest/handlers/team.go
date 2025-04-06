package handlers

import (
	"backend/internal/models"
	"backend/internal/transport/rest/req"
	"backend/internal/transport/rest/res"
	"backend/internal/transport/service"
	"backend/internal/utils"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

func (h *TeamHandler) GetAllRoles(c echo.Context) error {

	roles, err := h.team.GetAllRoles()
	if err != nil {
		code, msg := utils.InternalServerError("failed to get roles: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.GetAllRolesRes{
		Roles: roles,
	})
}
func (h *TeamHandler) GetRoleByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	// Преобразуем в UUID
	uuidID, err := uuid.Parse(id)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	role, err := h.team.GetRoleByID(uuidID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			code, msg := utils.NotFoundError()
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to get role: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, role)
}

func (h *TeamHandler) UpdateRole(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	var req req.UpdateRoleReq
	if err := c.Bind(&req); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	err = h.team.UpdateRole(uuidID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			code, msg := utils.NotFoundError()
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to update role: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.UpdateRoleRes{
		Message: "Role updated successfully",
	})
}

func (h *TeamHandler) DeleteRole(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	err = h.team.DeleteRole(uuidID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			code, msg := utils.NotFoundError()
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to delete role: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.DeleteRoleRes{
		Message: "Role deleted successfully",
	})
}

func (h *TeamHandler) CreateUser(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	if err := h.validate.Struct(user); err != nil {
		code, msg := utils.ValidationError()
		return c.JSON(code, msg)
	}
	userRole := c.Get("user_role").(string)

	err := h.team.CreateUser(&user, userRole)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			code, msg := utils.ConflictCustomError("User with this email already exists")
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to create user: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusCreated, res.CreateRoleRes{
		Message: "User created successfully",
	})
}

func (h *TeamHandler) GetAllUsers(c echo.Context) error {
	users, err := h.team.GetAllUsers()
	if err != nil {
		code, msg := utils.InternalServerError("failed to get users: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.GetAllUsersRes{
		Users: users,
	})
}
func (h *TeamHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	user, err := h.team.GetUserByID(uuidID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			code, msg := utils.NotFoundError()
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to get user: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *TeamHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	var req req.UpdateUserReq
	if err := c.Bind(&req); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}
	userRole := c.Get("user_role").(string)

	err = h.team.UpdateUser(uuidID, &req, userRole)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			code, msg := utils.NotFoundError()
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to update user: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.UpdateRoleRes{
		Message: "User updated successfully",
	})
}

func (h *TeamHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}
	userRole := c.Get("user_role").(string)

	uuidID, err := uuid.Parse(id)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	err = h.team.DeleteUser(uuidID, userRole)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			code, msg := utils.NotFoundError()
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to delete user: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.DeleteRoleRes{
		Message: "User deleted successfully",
	})
}

func (h *TeamHandler) UpdatePassword(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	uuidID, err := uuid.Parse(id)
	if err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	var req req.UpdatePasswordReq
	if err := c.Bind(&req); err != nil {
		code, msg := utils.BadRequestError()
		return c.JSON(code, msg)
	}

	err = h.team.UpdatePassword(uuidID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			code, msg := utils.NotFoundError()
			return c.JSON(code, msg)
		}
		code, msg := utils.InternalServerError("failed to update password: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, res.UpdateRoleRes{
		Message: "Password updated successfully",
	})
}
