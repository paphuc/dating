package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	TargetUserID primitive.ObjectID `json:"target_user_id" bson:"target_user_id"`
	Matched      bool               `json:"matched" bson:"matched"`
	CreateAt     time.Time          `json:"created_at" bson:"created_at"`
}

type MatchRequest struct {
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id" validate:"required"`
	TargetUserID primitive.ObjectID `json:"target_user_id" bson:"target_user_id" validate:"required"`
	Matched      bool               `json:"matched" bson:"matched"`
}

type MatchResponse struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	TargetUserID primitive.ObjectID `json:"target_user_id" bson:"target_user_id"`
	Matched      bool               `json:"matched" bson:"matched"`
	CreateAt     time.Time          `json:"created_at" bson:"created_at"`
	TargetUser   []UserResGetInfo   `json:"target_user" bson:"target_user"`
}
