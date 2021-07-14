package types

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Match struct {
	ID           bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	UserID       bson.ObjectId `json:"user_id" bson:"user_id"`
	TargetUserID bson.ObjectId `json:"target_user_id" bson:"target_user_id"`
	Matched      bool          `json:"matched" bson:"matched"`
	CreateAt     time.Time     `json:"created_at" bson:"created_at"`
}

type MatchRequest struct {
	UserID       bson.ObjectId `json:"user_id" bson:"user_id" validate:"required"`
	TargetUserID bson.ObjectId `json:"target_user_id" bson:"target_user_id" validate:"required"`
	Matched      bool          `json:"matched" bson:"matched"`
}

type MatchResponse struct {
	ID           bson.ObjectId    `json:"_id" bson:"_id,omitempty"`
	UserID       bson.ObjectId    `json:"user_id" bson:"user_id"`
	TargetUserID bson.ObjectId    `json:"target_user_id" bson:"target_user_id"`
	Match        bool             `json:"matched" bson:"matched"`
	CreateAt     time.Time        `json:"created_at" bson:"created_at"`
	TargetUser   []UserResGetInfo `json:"target_user" bson:"target_user"`
}
