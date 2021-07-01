package types

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID       bson.ObjectId   `json:"_id" bson:"_id,omitempty"`
	Name     string          `json:"name" bson:"name"`
	Email    string          `json:"email" bson:"email"`
	Password string          `json:"password" bson:"password"`
	Avatars  []string        `json:"avatar" bson:"avatars"` // arr path image
	Gender   string          `json:"gender" bson:"gender"`
	Country  string          `json:"country" bson:"country"`
	Hobby    []Hobby         `json:"hobby" bson:"hobby"`
	About    string          `json:"about" bson:"about"`
	LikeID   []bson.ObjectId `json:"like_id" bson:"like_id"`
	MatchID  []bson.ObjectId `json:"match_id" bson:"match_id"`
	CreateAt time.Time       `json:"created_at" bson:"created_at"`
	UpdateAt time.Time       `json:"updated_at" bson:"updated_at"`
}
type UserSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserFieldInToken struct {
	ID    bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
}
type UserResponseSignUp struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}
type UserLogin struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
