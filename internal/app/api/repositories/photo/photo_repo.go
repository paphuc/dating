package pictures

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

// this method helps insert picture for user
func (r *MongoRepository) Insert(ctx context.Context, pics types.Photo) error {
	s := r.session.Clone()
	defer s.Close()

	err := r.collection(s).Insert(pics)

	return err
}

// this method helps find pictures by id user
func (r *MongoRepository) FindByID(ctx context.Context, id bson.ObjectId) (*types.Photo, error) {
	s := r.session.Clone()
	defer s.Close()

	var pics *types.Photo
	err := r.collection(s).FindId(id).One(&pics)

	return pics, err
}

func (r *MongoRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("photos")
}
