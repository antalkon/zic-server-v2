package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/transport/rest/req"
	"fmt"

	"github.com/google/uuid"
)

type TeamService struct {
	teamRepo *repository.TeamRepository
}

func NewTeamService(repo *repository.TeamRepository) *TeamService {
	return &TeamService{teamRepo: repo}
}

func (s *TeamService) CreateRole(role *models.Role) error {
	exists, err := s.teamRepo.RoleExistsByName(role.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("role with name '%s' already exists", role.Name)
	}

	role.ID = uuid.New()
	return s.teamRepo.CreateRole(role)
}

func (s *TeamService) GetAllRoles() ([]models.Role, error) {
	roles, err := s.teamRepo.GetAllRoles()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *TeamService) GetRoleByID(id uuid.UUID) (*models.Role, error) {
	role, err := s.teamRepo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role with ID '%s' not found", id.String())
	}
	return role, nil
}

func (s *TeamService) UpdateRole(id uuid.UUID, req *req.UpdateRoleReq) error {
	roleInDB, err := s.teamRepo.GetRoleByID(id)
	if err != nil {
		return err
	}
	if roleInDB == nil {
		return fmt.Errorf("role with ID '%s' not found", id)
	}

	if req.Name != nil {
		roleInDB.Name = *req.Name
	}
	if req.Desc != nil {
		roleInDB.Desc = *req.Desc
	}

	return s.teamRepo.UpdateRole(roleInDB)
}

func (s *TeamService) DeleteRole(id uuid.UUID) error {
	role, err := s.teamRepo.GetRoleByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return fmt.Errorf("role with ID '%s' not found", id)
	}

	count, err := s.teamRepo.GetUsersCountByRoleID(id)
	if count > 0 {
		return fmt.Errorf("role with ID '%s' has users assigned", id)
	}
	return s.teamRepo.DeleteRole(id)

}
