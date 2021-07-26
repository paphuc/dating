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
type UserResGetInfoInRoom struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name" validate:"omitempty,max=60"`
	Avatar string             `json:"avatar" bson:"avatar"` // arr path media
	Gender string             `json:"gender" bson:"gender" validate:"omitempty,max=60"`
}
type MatchRoomResponse struct {
	ID          primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	User        []UserResGetInfoInRoom `json:"users" bson:"users"`
	LastMessage *Message               `json:"last_message" bson:"last_message,omitempty"`
	CreateAt    time.Time              `json:"created_at" bson:"created_at"`
}
