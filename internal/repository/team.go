package repository

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamnRepository(db *gorm.DB) *TeamRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateRole(role *models.Role) error {
	return r.db.Table("roles").Create(role).Error
}
func (r *TeamRepository) RoleExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Table("roles").Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *TeamRepository) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Table("roles").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
func (r *TeamRepository) GetRoleByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Table("roles").Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *TeamRepository) UpdateRole(role *models.Role) error {
	return r.db.Table("roles").Where("id = ?", role.ID).Updates(role).Error
}

func (r *TeamRepository) DeleteRole(id uuid.UUID) error {
	return r.db.Table("roles").Where("id = ?", id).Delete(&models.Role{}).Error
}
func (r *TeamRepository) GetUsersCountByRoleID(roleID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Table("users").Where("role_id = ?", roleID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (r *TeamRepository) UserExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Table("users").Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *TeamRepository) CreateUser(user *models.User) error {
	return r.db.Table("users").Create(user).Error
}
