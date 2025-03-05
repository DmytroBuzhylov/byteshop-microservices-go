package models

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"ID"`
	ProductID     string    `gorm:"not null;index" json:"productID"`
	BuyerID       string    `gorm:"not null;index" json:"buyerID"`
	SellerID      string    `gorm:"not null;index" json:"sellerID"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Fee           float64   `gorm:"not null" json:"fee"`
	NetAmount     float64   `gorm:"not null" json:"netAmount"`
	Status        string    `gorm:"type:varchar(50);not null" json:"status"`
	PaymentMethod string    `gorm:"type:varchar(50);not null" json:"paymentMethod"`
	TransactionID string    `gorm:"type:varchar(100);unique" json:"transactionID"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
