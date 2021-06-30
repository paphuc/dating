package types

import "github.com/globalsign/mgo/bson"

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
	Like_id  []bson.ObjectId `json:"like_id" bson:"like_id"`
	Math     []bson.ObjectId `json:"match" bson:"math"`
}
type UserSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
