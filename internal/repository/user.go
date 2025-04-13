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
	err := u.db.First(&user, id).Error
	if err != nil {
		return user, err
	}
	// Обнуляем хэш пароля, чтобы не возвращать
	user.PasswordHash = ""
	return user, nil
}
