package repository

import "gorm.io/gorm"

type TunelRepository struct {
	db *gorm.DB
}

func NewTunelRepository(db *gorm.DB) *TunelRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &TunelRepository{db: db}
}
