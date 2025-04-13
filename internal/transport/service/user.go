package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/transport/rest/req"
	"backend/pkg/hash"
	"fmt"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, data *req.UpdateUserDataReq) error {
	userUpdates := make(map[string]interface{})

	if data.FirstName != nil {
		userUpdates["name"] = *data.FirstName
	}
	if data.SurName != nil {
		userUpdates["surname"] = *data.SurName
	}
	if data.Email != nil {
		userUpdates["email"] = *data.Email
	}
	if data.Phone != nil {
		userUpdates["phone"] = *data.Phone
	}

	return s.userRepo.UpdateUserFields(id, userUpdates)
}

func (s *UserService) UpdatePassword(id uuid.UUID, data *req.UpdateDataPasswordReq) error {
	// Получаем текущий хеш из БД
	currentHash, err := s.userRepo.GetPasswordHash(id)
	if err != nil {
		return fmt.Errorf("failed to get current password hash: %w", err)
	}

	// Сравниваем старый пароль
	if err := hash.ComparePassword(data.OldPassword, currentHash); err != nil {
		return fmt.Errorf("old password is incorrect")
	}

	// Генерируем хеш нового пароля
	newHash, err := hash.GenerateHash(data.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to generate new password hash: %w", err)
	}

	// Обновляем хеш в БД
	if err := s.userRepo.UpdatePasswordHash(id, newHash); err != nil {
		return fmt.Errorf("failed to update password hash: %w", err)
	}

	return nil
}
