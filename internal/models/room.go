package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Number      int        `json:"number"`
	Description string     `gorm:"size:255" json:"description,omitempty"`
	ImageID     uuid.UUID  `gorm:"type:uuid" json:"image_id,omitempty"`
	Status      string     `gorm:"size:255;not null" json:"status,omitempty"` // active | blocked
	CreatedAt   time.Time  `gorm:"not null" json:"created_at"`
	Computers   []Computer `gorm:"foreignKey:Room" json:"computers,omitempty"`
}
