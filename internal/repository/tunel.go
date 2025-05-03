package repository

import (
	"backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TunelRepository struct {
	db *gorm.DB
}

func NewTunelRepository(db *gorm.DB) *TunelRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &TunelRepository{db: db}
}

func (r *TunelRepository) GetPcById(id uuid.UUID) (*models.Computer, error) {
	var pc models.Computer
	if err := r.db.First(&pc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pc, nil
}

func (r *TunelRepository) UpdateClientVersion(id uuid.UUID, version string) error {
	if err := r.db.Model(&models.Computer{}).Where("id = ?", id).Update("client_version", version).Error; err != nil {
		return err
	}
	return nil
}
func (r *TunelRepository) UpdateOS(id uuid.UUID, os string) error {
	if err := r.db.Model(&models.Computer{}).Where("id = ?", id).Update("os", os).Error; err != nil {
		return err
	}
	return nil
}
func (r *TunelRepository) UpdatePublicIP(id uuid.UUID, ip string) error {
	if err := r.db.Model(&models.Computer{}).Where("id = ?", id).Update("public_ip", ip).Error; err != nil {
		return err
	}
	return nil
}
func (r *TunelRepository) UpdateLocalIP(id uuid.UUID, ip string) error {
	if err := r.db.Model(&models.Computer{}).Where("id = ?", id).Update("local_ip", ip).Error; err != nil {
		return err
	}
	return nil
}
func (r *TunelRepository) UpdateLastActivity(id uuid.UUID) error {
	if err := r.db.Model(&models.Computer{}).Where("id = ?", id).Update("last_activity", gorm.Expr("NOW()")).Error; err != nil {
		return err
	}
	return nil
}

func (r *TunelRepository) ChangeStatus(id uuid.UUID, online bool) error {
	status := "off"
	if online {
		status = "on"
	}
	if err := r.db.Model(&models.Computer{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return err
	}
	return nil
}
