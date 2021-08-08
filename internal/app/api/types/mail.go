package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailVerification struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email,omitempty" validate:"omitempty,email"`
	Code        string             `json:"code" bson:"code,omitempty"`
	CreatedTime time.Time          `json:"created_time" bson:"created_time"`
	Verified    bool               `json:"verified" bson:"verified"`
}
