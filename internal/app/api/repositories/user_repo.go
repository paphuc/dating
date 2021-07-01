package user

import (
	"context"
	"dating/internal/app/api/types"
	"dating/internal/pkg/jwt"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

var (
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

//  This method helps to register new member
func (r *MongoRepository) SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error) {
	s := r.session.Clone()
	defer s.Close()
	var user types.User

	//check email exists
	if err := r.collection(s).Find(bson.M{"email": UserSignUp.Email}).One(&user); err == nil {
		return nil, errors.Wrap(errors.New("email email exits"), "email exits, can't insert user")
	}

	UserSignUp.Password, _ = jwt.HashPassword(UserSignUp.Password)
	user.CreateAt = time.Now()

	if err := r.collection(s).Insert(types.User{
		Name:     UserSignUp.Name,
		Email:    UserSignUp.Email,
		Password: UserSignUp.Password,
		CreateAt: time.Now()}); err != nil {
		return nil, errors.Wrap(err, "can't insert user")
	}

	error := r.collection(s).Find(bson.M{"email": UserSignUp.Email}).One(&user)
	if error != nil {
		return nil, errors.Wrap(error, "err insert user")
	}

	var tokenString string
	tokenString, err := jwt.GenToken(types.UserFieldInToken{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})

	if err != nil {
		return nil, errors.Wrap(err, "can't insert user")
	}

	return &types.UserResponseSignUp{
		Name:  UserSignUp.Name,
		Email: UserSignUp.Email,
		Token: tokenString}, nil
}

// This method helps user login
func (r *MongoRepository) Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error) {
	s := r.session.Clone()
	defer s.Close()
	var user types.User
	if err := r.collection(s).Find(bson.M{"email": UserLogin.Email}).One(&user); err != nil {
		return nil, errors.Wrap(errors.New("not found email exits"), "email not exists, can't find user")
	}
	// isCorrectPassword true
	if !jwt.IsCorrectPassword(UserLogin.Password, user.Password) {
		return nil, errors.Wrap(errors.New("password incorrect"), "password incorrect")
	}

	var tokenString string
	tokenString, err := jwt.GenToken(types.UserFieldInToken{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email})

	if err != nil {
		return nil, errors.Wrap(err, "can't insert user")
	}

	return &types.UserResponseSignUp{
		Name:  user.Name,
		Email: user.Email,
		Token: tokenString}, nil

}

func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.User, error) {
	s := r.session.Clone()
	defer s.Close()

	var user *types.User
	if err := r.collection(s).Find(bson.M{"_id": id}).One(&user); err != nil {
		return nil, errors.Wrap(err, "failed to find the given user from database")
	}

	return user, nil
}

func (r *MongoRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("users")
}
