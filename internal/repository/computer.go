package repository

import (
	"backend/internal/models"
	"backend/internal/transport/rest/req"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComputerRepository struct {
	db *gorm.DB
}

func NewComputerRepository(db *gorm.DB) *ComputerRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &ComputerRepository{db: db}
}

func (r *ComputerRepository) Create(computer *models.Computer) error {
	if err := r.db.Create(computer).Error; err != nil {
		return err
	}
	return nil
}

func (r *ComputerRepository) CheckRoomID(roomID string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Room{}).Where("id = ?", roomID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ComputerRepository) GetAll() ([]models.Computer, error) {
	var computers []models.Computer
	if err := r.db.Find(&computers).Error; err != nil {
		return nil, err
	}
	return computers, nil
}

func (r *ComputerRepository) GetByID(id string) (*models.Computer, error) {
	var computer models.Computer
	if err := r.db.First(&computer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &computer, nil
}

func (r *ComputerRepository) GetRoomComputersAllByRoomId(roomID uuid.UUID) ([]models.Computer, error) {
	var computers []models.Computer
	if err := r.db.Where("room = ?", roomID).Find(&computers).Error; err != nil {
		return nil, err
	}
	return computers, nil
}
func (r *ComputerRepository) UpdateComputer(id uuid.UUID, data *req.UpdateComputerReq) error {
	updates := map[string]interface{}{}

	if data.Name != nil {
		updates["name"] = *data.Name
	}
	if data.Location != nil {
		updates["location"] = *data.Location
	}
	if data.Description != nil {
		updates["description"] = *data.Description
	}
	if data.PublicIP != nil {
		updates["public_ip"] = *data.PublicIP
	}
	if data.LocalIP != nil {
		updates["local_ip"] = *data.LocalIP
	}
	if data.OS != nil {
		updates["os"] = *data.OS
	}
	if data.ClientVersion != nil {
		updates["client_version"] = *data.ClientVersion
	}
	if data.Comment != nil {
		updates["comment"] = *data.Comment
	}

	if len(updates) == 0 {
		return nil // ничего не обновляем
	}

	return r.db.Model(&models.Computer{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ComputerRepository) DeleteComputer(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Computer{}).Error
}
