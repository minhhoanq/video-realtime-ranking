package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InteractionType
type InteractionType string

const (
	SystemInteraction InteractionType = "SYSTEM"
	EmailInteraction  InteractionType = "EMAIL"
)

type Interaction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	VideoID         string             `bson:"video_id" json:"video_id"`
	InteractionType string             `bson:"interaction_type" json:"interaction_type"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type SendInteractionRequest struct {
	UserID          string `bson:"user_id" json:"user_id"`
	VideoID         string `bson:"video_id" json:"video_id"`
	InteractionType string `bson:"interaction_type" json:"interaction_type"`
}

type SendInteractionResponse struct {
	ID      string `bson:"_id,omitempty" json:"id"`
	UserID  string `bson:"user_id" json:"user_id"`
	VideoID string `bson:"video_id" json:"video_id"`
}
