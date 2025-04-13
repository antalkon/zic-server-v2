package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/transport/rest/req"
	"fmt"

	"github.com/google/uuid"
)

type RoomService struct {
	roomRepo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{roomRepo: repo}
}

func (s *RoomService) CreateRoom(room *models.Room, userRole string) (string, error) {
	if userRole != "admin" {
		return "", fmt.Errorf("only admins can create rooms")
	}
	uuid := uuid.New()
	room.ID = uuid
	room.Status = "active"
	roomID, err := s.roomRepo.CreateRoom(room)
	if err != nil {
		return "", err
	}
	return roomID, nil
}

func (s *RoomService) GetAllRooms() ([]models.Room, error) {
	rooms, err := s.roomRepo.GetAllRooms()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
func (s *RoomService) GetRoomByID(id string) (*models.Room, error) {
	room, err := s.roomRepo.GetRoomByID(id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (s *RoomService) UpdateRoom(roomID string, update *req.UpdateRoomRequest, userRole string) error {
	if userRole != "admin" {
		return fmt.Errorf("only admins can update rooms")
	}

	existingRoom, err := s.roomRepo.GetRoomByID(roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	// Обновляем только пришедшие поля
	if update.Name != nil {
		existingRoom.Name = *update.Name
	}
	if update.Number != nil {
		existingRoom.Number = *update.Number
	}
	if update.Description != nil {
		existingRoom.Description = *update.Description
	}
	if update.ImageId != nil {
		imageUUID, err := uuid.Parse(*update.ImageId)
		if err != nil {
			return fmt.Errorf("invalid image_id format: %w", err)
		}
		existingRoom.ImageID = imageUUID
	}

	return s.roomRepo.UpdateRoom(existingRoom)
}

func (s *RoomService) DeleteRoom(id string, userRole string) error {
	if userRole != "admin" {
		return fmt.Errorf("only admins can delete rooms")
	}
	err := s.roomRepo.DeleteRoom(id)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}
