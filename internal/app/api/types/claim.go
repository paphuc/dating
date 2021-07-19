package types

import (
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email string             `json:"email"`
	Name  string             `json:"name"`
	jwt.StandardClaims
}
