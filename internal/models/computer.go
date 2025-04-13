package models

import (
	"time"

	"github.com/google/uuid"
)

type Computer struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Room          uuid.UUID `gorm:"type:uuid;not null" json:"room_id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Location      string    `gorm:"size:255" json:"location,omitempty"`
	Description   string    `gorm:"size:255" json:"description,omitempty"`
	PublicIP      string    `gorm:"size:255;not null" json:"public_ip"`
	LocalIP       string    `gorm:"size:255" json:"local_ip,omitempty"`
	OS            string    `gorm:"size:255;not null" json:"os"`
	ClientVersion string    `gorm:"size:255;not null" json:"client_version"`
	Secret        string    `gorm:"size:255;not null" json:"secret"`
	Status        string    `gorm:"size:255;not null" json:"status"` // on | off
	Blocked       bool      `gorm:"not null;default:false" json:"blocked"`
	Bage          string    `gorm:"size:255;not null;default:'none'" json:"bage"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	LastActivity  time.Time `gorm:"not null" json:"last_activity"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	Comment       string    `gorm:"size:255" json:"comment,omitempty"`
}
