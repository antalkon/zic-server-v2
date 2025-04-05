package service

import (
	"backend/internal/models"
	"backend/internal/repository"
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
