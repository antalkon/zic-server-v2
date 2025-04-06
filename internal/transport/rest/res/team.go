package res

import "backend/internal/models"

type CreateRoleRes struct {
	Message string `json:"message"`
}

type GetAllRolesRes struct {
	Roles []models.Role `json:"roles"`
}

type UpdateRoleRes struct {
	Message string `json:"message"`
}

type DeleteRoleRes struct {
	Message string `json:"message"`
}

type GetAllUsersRes struct {
	Users []models.User `json:"users"`
}
