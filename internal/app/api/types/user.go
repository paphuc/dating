package types

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID           bson.ObjectId `json:"_id" bson:"_id,omitempty" validate:"required"`
	Name         string        `json:"name" bson:"name" validate:"required"`
	Email        string        `json:"email" bson:"email" validate:"omitempty,email"`
	Birthday     time.Time     `json:"birthday" bson:"birthday" validate:"required"`
	Relationship string        `json:"relationship" bson:"relationship" validate:"omitempty,max=60" `
	LookingFor   string        `json:"looking_for" bson:"looking_for" validate:"omitempty,max=60"`
	Password     string        `json:"password" bson:"password"`
	Media        []string      `json:"media" bson:"media"` // arr path media
	Gender       string        `json:"gender" bson:"gender" validate:"required,eq=Male|eq=Female|eq=Both"`
	Sex          string        `json:"sex" bson:"sex" validate:"omitempty,max=60"`
	Country      string        `json:"country" bson:"country" validate:"required,max=60"`
	Hobby        []string      `json:"hobby" bson:"hobby"`
	Disable      bool          `json:"disable" bson:"disable"`
	About        string        `json:"about" bson:"about" validate:"omitempty,max=256"`
	CreateAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdateAt     time.Time     `json:"updated_at" bson:"updated_at"`
}

type UserResGetInfo struct {
	ID           bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Name         string        `json:"name" bson:"name" validate:"omitempty,max=60"`
	Email        string        `json:"email" bson:"email"`
	Birthday     time.Time     `json:"birthday" bson:"birthday"`
	Relationship string        `json:"relationship" bson:"relationship" validate:"omitempty,max=60"`
	LookingFor   string        `json:"looking_for" bson:"looking_for" validate:"omitempty,max=60"`
	Media        []string      `json:"media" bson:"media"` // arr path media
	Gender       string        `json:"gender" bson:"gender" validate:"omitempty,max=60"`
	Sex          string        `json:"sex" bson:"sex" validate:"omitempty,max=60"`
	Country      string        `json:"country" bson:"country" validate:"omitempty,max=60"`
	Hobby        []string      `json:"hobby" bson:"hobby"`
	About        string        `json:"about" bson:"about" validate:"omitempty,max=256"`
	CreateAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdateAt     time.Time     `json:"updated_at" bson:"updated_at"`
}

type UserSignUp struct {
	Name     string `json:"name" validate:"required,max=60"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
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

type GetListUsersResponse struct {
	Pagination
	ListUsers []*UserResGetInfo `json:"listUsers"`
}
type UserLogin struct {
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type DisableBody struct {
	Disable *bool `json:"disable" bson:"disable" validate:"required"`
}
