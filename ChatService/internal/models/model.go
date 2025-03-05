package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Sender    string    `bson:"sender" json:"sender"`
	Content   string    `bson:"content" json:"content"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Chat struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Participants []string           `bson:"participants" json:"participants"`
	Messages     []Message          `bson:"messages" json:"messages"`
	LastMessage  *Message           `bson:"lastMessage,omitempty" json:"lastMessage,omitempty"`
}

type Review struct {
	Name      string    `json:"name" bson:"name"`
	Content   string    `json:"content" bson:"content"`
	Grade     int       `json:"grade" bson:"grade"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type ProductReviews struct {
	ProductID string   `json:"productID" bson:"productID"`
	Reviews   []Review `json:"reviews" bson:"reviews"`
}
