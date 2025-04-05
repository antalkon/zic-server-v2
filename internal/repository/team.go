package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamnRepository(db *gorm.DB) *TeamRepository {
	if db == nil {
		panic("Database connection is nil in repository")
	}
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateRole(role *models.Role) error {
	return r.db.Table("roles").Create(role).Error
}
func (r *TeamRepository) RoleExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Table("roles").Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
