package match

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

// this method helps insert match
func (r *MongoRepository) Insert(ctx context.Context, match types.Match) error {
	s := r.session.Clone()
	defer s.Close()

	err := r.collection(s).Insert(match)

	return err
}

// This method helps insert match
func (r *MongoRepository) DeleteMatch(ctx context.Context, id string) error {
	s := r.session.Clone()
	defer s.Close()

	err := r.collection(s).RemoveId(bson.ObjectIdHex(id))

	return err
}

// This method helps get basic info match by id
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.Match, error) {
	s := r.session.Clone()
	defer s.Close()

	var match *types.Match
	err := r.collection(s).FindId(bson.ObjectIdHex(id)).One(&match)

	return match, err
}

// This method help check A vs B by Match
func (r *MongoRepository) CheckAB(ctx context.Context, idUser, idTargetUser string, matched bool) (*types.Match, error) {
	s := r.session.Clone()
	defer s.Close()

	filter := bson.M{
		"user_id":        bson.ObjectIdHex(idUser),
		"target_user_id": bson.ObjectIdHex(idTargetUser),
		"matched":        matched,
	}

	var match *types.Match
	err := r.collection(s).Find(filter).One(&match)

	return match, err
}

// this method help get record when user A liked user B
func (r *MongoRepository) FindALikeB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error) {
	s := r.session.Clone()
	defer s.Close()

	filter := bson.M{
		"user_id":        bson.ObjectIdHex(idUser),
		"target_user_id": bson.ObjectIdHex(idTargetUser),
	}

	var match *types.Match
	err := r.collection(s).Find(filter).One(&match)

	return match, err
}
func (r *MongoRepository) FindAMatchB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error) {
	s := r.session.Clone()
	defer s.Close()

	filter := bson.M{
		"$or": []interface{}{
			bson.M{
				"user_id":        bson.ObjectIdHex(idUser),
				"target_user_id": bson.ObjectIdHex(idTargetUser),
			},
			bson.M{
				"user_id":        bson.ObjectIdHex(idTargetUser),
				"target_user_id": bson.ObjectIdHex(idUser),
			},
		},
		"matched": true,
	}
	var match *types.Match
	err := r.collection(s).Find(filter).One(&match)

	return match, err
}

// this method help get update match true when A,B liked
func (r *MongoRepository) UpdateMatchByID(ctx context.Context, id string) error {
	s := r.session.Clone()
	defer s.Close()

	updatedUser := bson.M{"$set": bson.M{"matched": true}}
	err := r.collection(s).UpdateId(bson.ObjectIdHex(id), updatedUser)

	return err
}

// this method help get Upsert match
func (r *MongoRepository) UpsertMatch(ctx context.Context, match types.Match) error {
	s := r.session.Clone()
	defer s.Close()
	filter := bson.M{
		"user_id":        match.UserID,
		"target_user_id": match.TargetUserID,
	}
	updatedMath := bson.M{"$set": bson.M{
		"user_id":        match.UserID,
		"target_user_id": match.TargetUserID,
		"matched":        false,
		"created_at":     time.Now(),
	}}
	_, err := r.collection(s).Upsert(filter, updatedMath)
	return err
}

// this method help get list like
func (r *MongoRepository) GetListLiked(ctx context.Context, idUser string) ([]*types.Match, error) {
	s := r.session.Clone()
	defer s.Close()

	filter := bson.M{
		"user_id": bson.ObjectIdHex(idUser),
		"matched": false,
	}
	var match []*types.Match
	err := r.collection(s).Find(filter).All(&match)

	return match, err
}

// this method help get list matched
func (r *MongoRepository) GetListMatched(ctx context.Context, idUser string) ([]*types.Match, error) {
	s := r.session.Clone()
	defer s.Close()

	filter := bson.M{
		"$or": []interface{}{
			bson.M{"user_id": bson.ObjectIdHex(idUser)},
			bson.M{"target_user_id": bson.ObjectIdHex(idUser)},
		},
		"matched": true,
	}
	var match []*types.Match
	err := r.collection(s).Find(filter).All(&match)

	return match, err
}

// this method help get list matched include info
func (r *MongoRepository) GetListMatchedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error) {
	s := r.session.Clone()
	defer s.Close()

	queryTest := []bson.M{
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
	}
	var listMatched []*types.UserResGetInfo
	err := r.collection(s).Pipe(queryTest).All(&listMatched)

	return listMatched, err
}

// this method help get list matched include info
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
	}

	var listMatched []*types.UserResGetInfo

	err := r.collection(s).Pipe(query).All(&listMatched)
	return listMatched, err
}

func (r *MongoRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("matches")
}
