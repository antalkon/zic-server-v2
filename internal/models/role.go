package models

import (
	"github.com/google/uuid"
)

type Role struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name  string    `gorm:"size:255;not null" json:"name"`
	Desc  string    `gorm:"size:255" json:"desc"`
	Users []User    `gorm:"foreignKey:RoleID"`
}
