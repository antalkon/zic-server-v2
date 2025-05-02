package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/transport/rest/req"
	tokenjwt "backend/pkg/token_jwt"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

type ComputerService struct {
	computerRepo *repository.ComputerRepository
}

func NewComputerService(repo *repository.ComputerRepository) *ComputerService {
	return &ComputerService{computerRepo: repo}
}

func (s *ComputerService) CreateComputer(computer *models.Computer, userRole string) (string, error) {
	if userRole != "admin" {
		return "", fmt.Errorf("user does not have permission to create a computer")
	}

	// Проверка существования комнаты
	roomExists, err := s.computerRepo.CheckRoomID(computer.Room.String())
	if err != nil {
		return "", fmt.Errorf("failed to check room ID: %w", err)
	}
	if !roomExists {
		return "", fmt.Errorf("room ID does not exist")
	}

	// Заполнение полей
	computer.ID = uuid.New()
	computer.Status = "off"
	computer.Blocked = false
	computer.CreatedAt = time.Now()
	computer.LastActivity = time.Now()
	computer.UpdatedAt = time.Now()
	computer.Secret = generateSecret(12)

	// Генерация токена
	jwt, err := tokenjwt.GenerateTunnelToken(computer.ID.String())
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Сохраняем компьютер в БД
	if err := s.computerRepo.Create(computer); err != nil {
		return "", fmt.Errorf("failed to create computer: %w", err)
	}

	return jwt, nil
}

func generateSecret(length int) string {
	var charset = "abcdefghijklmnopqrstuvwxyz"
	rand.Seed(uint64(time.Now().UnixNano()))
	secret := make([]byte, length)
	for i := range secret {
		secret[i] = charset[rand.Intn(len(charset))]
	}
	return string(secret)
}

func (s *ComputerService) GetAllComputers() ([]models.Computer, error) {
	computers, err := s.computerRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get computers: %w", err)
	}
	return computers, nil
}

func (s *ComputerService) GetComputerByID(id string) (*models.Computer, error) {
	computer, err := s.computerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get computer by ID: %w", err)
	}
	return computer, nil
}

func (s *ComputerService) GetRoomComputersAll(roomID uuid.UUID) ([]models.Computer, error) {
	computers, err := s.computerRepo.GetRoomComputersAllByRoomId(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get computers: %w", err)
	}
	return computers, nil
}

func (s *ComputerService) UpdateComputer(id uuid.UUID, data *req.UpdateComputerReq, userRoke string) error {
	if userRoke != "admin" {
		return fmt.Errorf("user does not have permission to update a computer")
	}
	return s.computerRepo.UpdateComputer(id, data)
}

func (s *ComputerService) DeleteComputer(id uuid.UUID, userRole string) error {
	if userRole != "admin" {
		return fmt.Errorf("user does not have permission to delete a computer")
	}
	return s.computerRepo.DeleteComputer(id)
}
