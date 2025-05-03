package repository

import (
	"backend/internal/models"
	"backend/internal/transport/rest/res"

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

func (r *RoomRepository) GetRoomComputers(roomID string) ([]models.Computer, error) {
	var computers []models.Computer
	err := r.db.Where("room = ?", roomID).Find(&computers).Error
	if err != nil {
		return nil, err
	}
	return computers, nil
}

func (r *RoomRepository) GetRoomStatusByID(roomID string) (*res.RoomStatusRes, error) {
	var (
		total   int64
		online  int64
		offline int64
		blocked int64
	)

	// Общее количество ПК
	if err := r.db.Model(&models.Computer{}).Where("room = ?", roomID).Count(&total).Error; err != nil {
		return nil, err
	}

	// Онлайн
	if err := r.db.Model(&models.Computer{}).Where("room = ? AND status = ?", roomID, "on").Count(&online).Error; err != nil {
		return nil, err
	}

	// Офлайн
	if err := r.db.Model(&models.Computer{}).Where("room = ? AND status = ?", roomID, "off").Count(&offline).Error; err != nil {
		return nil, err
	}

	// Заблокировано
	if err := r.db.Model(&models.Computer{}).Where("room = ? AND blocked = true", roomID).Count(&blocked).Error; err != nil {
		return nil, err
	}

	res := &res.RoomStatusRes{
		Status:    "active", // можно потом тянуть из самой комнаты
		PcTotal:   int(total),
		PcOnline:  int(online),
		PcOffline: int(offline),
		Blocked:   int(blocked),
		Free:      int(total - blocked),
	}
	return res, nil
}
