package repository

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserByID(id uuid.UUID) (models.User, error) {
	var user models.User
	err := u.db.Preload("Role").First(&user, "id = ?", id).Error
	if err != nil {
		return user, err
	}
	user.PasswordHash = ""
	return user, nil
}
func (u *UserRepository) GetPasswordHash(id uuid.UUID) (string, error) {
	var user models.User
	err := u.db.Select("password_hash").First(&user, "id = ?", id).Error
	if err != nil {
		return "", err
	}
	return user.PasswordHash, nil
}
func (u *UserRepository) UpdatePasswordHash(id uuid.UUID, hash string) error {
	return u.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("password_hash", hash).Error
}

func (u *UserRepository) UpdateUserFields(id uuid.UUID, fields map[string]interface{}) error {
	return u.db.Model(&models.User{}).Where("id = ?", id).Updates(fields).Error
}
