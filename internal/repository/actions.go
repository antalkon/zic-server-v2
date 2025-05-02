package repository

import "gorm.io/gorm"

type ActionsRepository struct {
	db *gorm.DB
}

func NewActionsRepository(db *gorm.DB) *ActionsRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &ActionsRepository{db: db}
}
