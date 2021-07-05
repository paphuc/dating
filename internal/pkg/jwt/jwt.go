package jwt

import (
	"dating/internal/app/api/types"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey, _ = base64.URLEncoding.DecodeString("dating21")
)

//Generate token for login or sign up
func GenToken(user types.UserFieldInToken) (string, error) {
	expirationTime := time.Now().Add(120 * time.Hour)
	claims := &types.Claims{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
func IsAuthorized(tokenpath string) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenpath, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Can't authorized token")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, errors.New("Can't authorized token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Can't authorized token")
}

// method hash password
func HashPassword(password string) (string, error) {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword), nil
}

// method compare Hash and Password
func IsCorrectPassword(password, hashedPasswordStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordStr), []byte(password))
	return err == nil
}
