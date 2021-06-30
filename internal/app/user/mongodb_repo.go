package user

import (
	"context"
	"dating/internal/app/types"
	"encoding/base64"
	"fmt"

	"time"

	"github.com/dgrijalva/jwt-go"
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey, _   = base64.URLEncoding.DecodeString("dating21")
	ErrNotFound = errors.New("not found")
)

type MongoRepository struct {
	session *mgo.Session
}

func NewMongoRepository(s *mgo.Session) *MongoRepository {
	return &MongoRepository{
		session: s,
	}
}

// Sign up method
func (r *MongoRepository) SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error) {
	s := r.session.Clone()
	defer s.Close()
	var User types.User

	//check email exists
	if err := r.collection(s).Find(bson.M{"email": UserSignUp.Email}).One(&User); err == nil {
		return nil, errors.Wrap(errors.New("email email exits"), "email exits, can't insert user")
	}

	UserSignUp.Password, _ = hashPassword(UserSignUp.Password)
	if err := r.collection(s).Insert(UserSignUp); err != nil {
		return nil, errors.Wrap(err, "can't insert user")
	}

	var tokenString string
	tokenString, err := GenToken(UserSignUp)

	if err != nil {
		return nil, errors.Wrap(err, "can't insert user")
	}

	return &types.UserResponseSignUp{
		Name:  UserSignUp.Name,
		Email: UserSignUp.Email,
		Token: tokenString}, nil
}

func (r *MongoRepository) Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error) {
	s := r.session.Clone()
	defer s.Close()
	var User types.User
	if err := r.collection(s).Find(bson.M{"email": UserLogin.Email}).One(&User); err != nil {
		return nil, errors.Wrap(errors.New("not found email exits"), "email not exists, can't find user")
	}
	// isCorrectPassword true
	isCorrectPassword := isCorrectPassword(UserLogin.Password, User.Password)

	if !isCorrectPassword {
		return nil, errors.Wrap(errors.New("password incorrect"), "password incorrect")
	}
	var tokenString string
	tokenString, err := GenToken(types.UserSignUp{Name: User.Name, Email: User.Email})

	if err != nil {
		return nil, errors.Wrap(err, "can't insert user")
	}

	return &types.UserResponseSignUp{
		Name:  User.Name,
		Email: User.Email,
		Token: tokenString}, nil

}

func (r *MongoRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("users")
}

//Generate token
func GenToken(user types.UserSignUp) (string, error) {
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

func hashPassword(password string) (string, error) {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))
	return string(hashedPassword), nil
}

func isCorrectPassword(password, hashedPasswordStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordStr), []byte(password))
	return err == nil
}
