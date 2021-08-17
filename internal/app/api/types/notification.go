package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	TokenDevice string             `json:"token_device" bson:"token_device"`
	CreateAt    time.Time          `json:"created_at" bson:"created_at"`
}
