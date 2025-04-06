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

func (r *TeamRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Role").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (r *TeamRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *TeamRepository) UpdateUser(user *models.User) error {
	return r.db.Table("users").Where("id = ?", user.ID).Updates(user).Error
}
func (r *TeamRepository) DeleteUser(id uuid.UUID) error {
	// Удаляем сначала токены
	if err := r.db.Where("user_id = ?", id).Delete(&models.RefreshToken{}).Error; err != nil {
		return err
	}

	// Потом сам юзер
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}
