package user

import (
	"context"
	"time"

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
	err := r.collection(s).Find(bson.M{"email": email, "disable": false}).One(&user)

	return user, err
}

// This method helps get basic info user
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.UserResGetInfo, error) {
	s := r.session.Clone()
	defer s.Close()

	var user *types.UserResGetInfo
	err := r.collection(s).Find(bson.M{"_id": bson.ObjectIdHex(id), "disable": false}).One(&user)

	return user, err
}

//  This method helps update info user
func (r *MongoRepository) UpdateUserByID(ctx context.Context, user types.User) error {
	s := r.session.Clone()
	defer s.Close()

	updatedUser := bson.M{"$set": bson.M{
		"name":         user.Name,
		"birthday":     user.Birthday,
		"relationship": user.Relationship,
		"looking_for":  user.LookingFor,
		"media":        user.Media,
		"gender":       user.Gender,
		"country":      user.Country,
		"hobby":        user.Hobby,
		"sex":          user.Sex,
		"about":        user.About,
		"updated_at":   time.Now(),
	}}
	err := r.collection(s).UpdateId(user.ID, updatedUser)

	return err
}

// This method helps Enable/Disable account
func (r *MongoRepository) DisableUserByID(ctx context.Context, idUser string, disable bool) error {
	s := r.session.Clone()
	defer s.Close()

	disableUpdate := bson.M{"$set": bson.M{
		"disable": disable,
	}}
	return r.collection(s).UpdateId(bson.ObjectIdHex(idUser), disableUpdate)
}

// This method helps get all user by page
func (r *MongoRepository) GetListUsers(ctx context.Context, ps types.PagingNSorting) ([]*types.UserResGetInfo, error) {
	s := r.session.Clone()
	defer s.Close()

	var result []*types.UserResGetInfo
	err := r.collection(s).Find(bson.M{
		"disable": false,
	}).Skip((ps.Page - 1) * ps.Size).Limit(ps.Size).All(&result)

	return result, err
}

// This method helps count number users in collection
func (r *MongoRepository) CountUser(ctx context.Context) (int, error) {
	s := r.session.Clone()
	defer s.Close()

	number, err := r.collection(s).Find(bson.M{
		"disable": false,
	}).Count()

	return number, err
}

// this method help get list matched include info
func (r *MongoRepository) GetListMatchedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error) {
	s := r.session.Clone()
	defer s.Close()

	query := []bson.M{
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"user_id": bson.ObjectIdHex(idUser)},
				bson.M{"target_user_id": bson.ObjectIdHex(idUser)},
			},
			"matched": true,
		}},
		{"$project": bson.M{
			"targer_id": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$user_id", bson.ObjectIdHex(idUser)}},
					"$target_user_id", "$user_id"},
			},
		},
		},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "targer_id",
			"foreignField": "_id",
			"as":           "target_user",
		}},
		{"$unwind": "$target_user"},
		{"$replaceRoot": bson.M{"newRoot": "$target_user"}},
		{"$match": bson.M{
			"disable": false,
		}},
	}
	var listMatched []*types.UserResGetInfo
	err := s.DB("").C("matches").Pipe(query).All(&listMatched)

	return listMatched, err
}

// this method help get list liked include info
func (r *MongoRepository) GetListlikedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error) {
	s := r.session.Clone()
	defer s.Close()

	filter := bson.M{
		"user_id": bson.ObjectIdHex(idUser),
		"matched": false,
	}
	query := []bson.M{
		{"$match": filter},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "target_user_id",
			"foreignField": "_id",
			"as":           "target_user",
		}},
		{"$unwind": "$target_user"},
		{"$replaceRoot": bson.M{"newRoot": "$target_user"}},
		{"$match": bson.M{
			"disable": false,
		}},
	}

	var listMatched []*types.UserResGetInfo

	err := s.DB("").C("matches").Pipe(query).All(&listMatched)
	return listMatched, err
}

func (r *MongoRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("users")
}
