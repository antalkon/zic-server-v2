package repository

import "gorm.io/gorm"

type FristRepository struct {
	db *gorm.DB
}

func NewFristRepository(db *gorm.DB) *FristRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &FristRepository{db: db}
}
