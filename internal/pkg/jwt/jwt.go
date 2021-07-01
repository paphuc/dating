package jwt

import (
	"dating/internal/app/api/types"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey, _ = base64.URLEncoding.DecodeString("dating21")
)

//Generate token for login or sign up
func GenToken(user types.UserFieldInToken) (string, error) {
	expirationTime := time.Now().Add(120 * time.Minute)
	claims := &types.Claims{
		Email: user.Email,
		Name:  user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func HashPassword(password string) (string, error) {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword), nil
}

func IsCorrectPassword(password, hashedPasswordStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordStr), []byte(password))
	return err == nil
}
