package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &RoomRepository{db: db}
}

func (r *RoomRepository) CreateRoom(room *models.Room) (string, error) {

	err := r.db.Create(room).Error
	if err != nil {
		return "", err
	}
	return room.ID.String(), nil
}

func (r *RoomRepository) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
func (r *RoomRepository) GetRoomByID(id string) (*models.Room, error) {
	var room models.Room
	err := r.db.Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) UpdateRoom(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *RoomRepository) DeleteRoom(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Room{}).Error
}
