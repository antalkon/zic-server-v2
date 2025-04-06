package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/transport/rest/req"
	"backend/pkg/hash"
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

func (s *TeamService) CreateUser(user *models.User, role string) error {
	if role != "admin" {
		return fmt.Errorf("only admin can create users")
	}
	exists, err := s.teamRepo.UserExistsByEmail(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user with email '%s' already exists", user.Email)
	}
	passwordHash, err := hash.GenerateHash(user.Password)
	if err != nil {
		return err
	}

	user.ID = uuid.New()
	user.PasswordHash = passwordHash

	return s.teamRepo.CreateUser(user)
}

func (s *TeamService) GetAllUsers() ([]models.User, error) {
	users, err := s.teamRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *TeamService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.teamRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user with ID '%s' not found", id.String())
	}
	return user, nil
}

func (s *TeamService) UpdateUser(id uuid.UUID, req *req.UpdateUserReq, role string) error {
	userInDB, err := s.teamRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if userInDB == nil {
		return fmt.Errorf("user with ID '%s' not found", id)
	}
	if role != "admin" {
		return fmt.Errorf("only admin can create users")
	}
	if req.Name != nil {
		userInDB.Name = *req.Name
	}
	if req.Surname != nil {
		userInDB.Surname = *req.Surname
	}
	if req.Email != nil {
		userInDB.Email = *req.Email
	}
	if req.Phone != nil {
		userInDB.Phone = *req.Phone
	}
	if req.Role != nil {
		role, err := s.teamRepo.GetRoleByID(uuid.MustParse(*req.Role))
		if err != nil {
			return err
		}
		if role == nil {
			return fmt.Errorf("role with ID '%s' not found", *req.Role)
		}
		userInDB.RoleID = role.ID
	}

	return s.teamRepo.UpdateUser(userInDB)
}

func (s *TeamService) DeleteUser(id uuid.UUID, role string) error {
	user, err := s.teamRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with ID '%s' not found", id)
	}
	if role != "admin" {
		return fmt.Errorf("only admin can create users")
	}

	return s.teamRepo.DeleteUser(id)
}

func (s *TeamService) UpdatePassword(id uuid.UUID, req *req.UpdatePasswordReq) error {
	userInDB, err := s.teamRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if userInDB == nil {
		return fmt.Errorf("user with ID '%s' not found", id)
	}

	passwordHash, err := hash.GenerateHash(req.NewPassword)
	if err != nil {
		return err
	}

	userInDB.PasswordHash = passwordHash

	return s.teamRepo.UpdateUser(userInDB)
}
