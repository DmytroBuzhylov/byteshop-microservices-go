package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"productID"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Type        string    `gorm:"not null" json:"type"`
	Link        string    `json:"link"`
	ImageURL    string    `gorm:"not null" json:"imageURL"`
	Category    string    `gorm:"not null" json:"category"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ProductItems struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductID uuid.UUID  `gorm:"not null" json:"productID"`
	Value     string     `gorm:"not null" json:"value" gorm:"not null"`
	IsSolid   bool       `json:"is_solid" gorm:"default:false"`
	SoldAt    *time.Time `json:"sold_at"`
}
