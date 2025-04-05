package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/internal/models"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(client *models.User) error {
	return r.db.Table("users").Create(client).Error
}

func (r *AuthRepository) GetUserByEmailNumber(email string) (*models.User, error) {
	var client models.User
	err := r.db.Table("users").
		Where("email = ?", email).
		First(&client).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &client, nil
}

func (r *AuthRepository) CheckUsersCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *AuthRepository) GetAdminRoleId() (uuid.UUID, error) {
	var role models.Role
	err := r.db.Table("roles").Where("name = ?", "admin").First(&role).Error
	if err != nil {
		return uuid.Nil, err
	}
	return role.ID, nil
}

func (r *AuthRepository) SaveRefreshToken(client *models.RefreshToken) error {
	return r.db.Table("refresh_tokens").Create(client).Error
}

func (r *AuthRepository) GetRefreshTokenById(id uuid.UUID) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.Table("refresh_tokens").Where("id = ?", id).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) GetUserById(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Table("users").Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) DeleteRefreshToken(tokenId uuid.UUID) error {
	return r.db.Where("id = ?", tokenId).Delete(&models.RefreshToken{}).Error
}
