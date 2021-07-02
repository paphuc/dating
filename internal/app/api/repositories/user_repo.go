package user

import (
	"context"
	"dating/internal/app/api/types"

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

// this method helps insert user
func (r *MongoRepository) Insert(ctx context.Context, user types.User) error {
	s := r.session.Clone()
	defer s.Close()

	err := r.collection(s).Insert(user)

	return err
}

// this method helps get user with email
func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	s := r.session.Clone()
	defer s.Close()

	var user *types.User
	err := r.collection(s).Find(bson.M{"email": email}).One(&user)

	return user, err
}

// This method helps get basic info user
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.UserResGetInfo, error) {
	s := r.session.Clone()
	defer s.Close()

	var user *types.UserResGetInfo
	err := r.collection(s).Find(bson.M{"_id": id}).One(&user)

	return user, err
}

func (r *MongoRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("users")
}
