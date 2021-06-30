package member

import (
	"context"

	"dating/internal/app/types"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// MongoRepository is MongoDB implementation of repository
type MongoRepository struct {
	session *mgo.Session
}

// NewMongoRepository return new MongoDB repository
func NewMongoRepository(s *mgo.Session) *MongoRepository {
	return &MongoRepository{
		session: s,
	}
}

// FindByID return member base on given id
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.Member, error) {
	s := r.session.Clone()
	defer s.Close()
	var member *types.Member
	if err := r.collection(s).Find(bson.M{"id": id}).One(&member); err != nil {
		return nil, errors.Wrap(err, "failed to find the given member from database")
	}

	return member, nil
}

func (r *MongoRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("members")
}
