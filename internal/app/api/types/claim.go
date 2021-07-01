package types

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
)

type Claims struct {
	ID    bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Email string        `json:"email"`
	Name  string        `json:"name"`
	jwt.StandardClaims
}
