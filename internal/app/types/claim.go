package types

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}
