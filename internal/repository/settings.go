package repository

import "gorm.io/gorm"

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &SettingsRepository{db: db}
}
