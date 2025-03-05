package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BuyerID   uuid.UUID `gorm:"not null" json:"buyerID"`
	ProductID uuid.UUID `gorm:"not null" json:"productID"`
	ItemID    string    `gorm:"not null" json:"itemID"`
	Amount    float64   `gorm:"not null" json:"amount"`
	Status    string    `gorm:"size:50;not null;default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
