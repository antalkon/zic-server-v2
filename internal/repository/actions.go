package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type ActionsRepository struct {
	db *gorm.DB
}

func NewActionsRepository(db *gorm.DB) *ActionsRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &ActionsRepository{db: db}
}

func (r *ActionsRepository) BlockComputer(computerID string) error {
	return r.db.Model(&models.Computer{}).
		Where("id = ?", computerID).
		Update("blocked", true).
		Error
}

func (r *ActionsRepository) UnblockComputer(computerID string) error {
	return r.db.Model(&models.Computer{}).
		Where("id = ?", computerID).
		Update("blocked", false).
		Error
}
