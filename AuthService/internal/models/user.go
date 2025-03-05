package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" `
	Name         string    `gorm:"not null" json:"name"`
	Role         string    `gorm:"not null;default:'buyer'" json:"role"`
	IsBanned     bool      `gorm:"default:false" json:"isBanned"`
	PayPalEmail  string    `gorm:"size:255" json:"payPalEmail"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}
