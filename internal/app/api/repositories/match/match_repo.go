package match

import (
	"context"
	"time"

	"dating/internal/app/api/types"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound = errors.New("not found")
)

type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(c *mongo.Client) *MongoRepository {
	return &MongoRepository{
		client: c,
	}
}

// this method helps insert match
func (r *MongoRepository) Insert(ctx context.Context, match types.Match) error {
	_, err := r.collection().InsertOne(context.TODO(), match)
	return err
}

// This method helps insert match
func (r *MongoRepository) DeleteMatch(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection().DeleteOne(context.TODO(), objectID)
	return err
}

// This method helps get basic info match by id
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.Match, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var match *types.Match
	err = r.collection().FindOne(ctx, types.Match{ID: objectID}).Decode(&match)
	return match, err
}

// This method help check A vs B by Match
func (r *MongoRepository) CheckAB(ctx context.Context, idUser, idTargetUser string, matched bool) (*types.Match, error) {
	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return nil, err
	}
	targetUserID, err := primitive.ObjectIDFromHex(idTargetUser)
	if err != nil {
		return nil, err
	}

	var match *types.Match
	err = r.collection().FindOne(ctx, types.Match{UserID: userID, TargetUserID: targetUserID, Matched: matched}).Decode(&match)

	return match, err
}

// this method help get record when user A liked user B
func (r *MongoRepository) FindALikeB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error) {
	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return nil, err
	}
	targetUserID, err := primitive.ObjectIDFromHex(idTargetUser)
	if err != nil {
		return nil, err
	}

	var match *types.Match
	err = r.collection().FindOne(ctx, types.Match{UserID: userID, TargetUserID: targetUserID}).Decode(&match)
	return match, err
}
func (r *MongoRepository) FindAMatchB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error) {
	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return nil, err
	}
	targetUserID, err := primitive.ObjectIDFromHex(idTargetUser)
	if err != nil {
		return nil, err
	}

	var match *types.Match
	if err := r.collection().FindOne(ctx, types.Match{UserID: userID, TargetUserID: targetUserID, Matched: true}).Decode(&match); err == nil {
		return match, err
	}

	if match == nil {
		if err := r.collection().FindOne(ctx, types.Match{UserID: targetUserID, TargetUserID: userID, Matched: true}).Decode(&match); err == nil {
			return match, err
		}
	}

	return match, err
}

// this method help get update match true when A,B liked
func (r *MongoRepository) UpdateMatchByID(ctx context.Context, id string) error {
	matchID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.D{
		{"$set", bson.D{
			{"matched", true},
		}},
	}

	_, err = r.collection().UpdateOne(ctx, types.Match{ID: matchID}, update)
	return err
}

// this method help get Upsert match
func (r *MongoRepository) UpsertMatch(ctx context.Context, match types.Match) error {
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

	opts := options.Update().SetUpsert(true)
	_, err := r.collection().UpdateOne(ctx, filter, updatedMath, opts)
	return err
}

// this method help get list like
func (r *MongoRepository) GetListLiked(ctx context.Context, idUser string) ([]*types.Match, error) {
	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return nil, err
	}

	var match []*types.Match
	cursor, err := r.collection().Find(ctx, types.Match{UserID: userID, Matched: true})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &match); err != nil {
		return nil, err
	}
	return match, err
}

// this method help get list matched
func (r *MongoRepository) GetListMatched(ctx context.Context, idUser string) ([]*types.Match, error) {
	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"$or": []interface{}{
			bson.M{"user_id": userID},
			bson.M{"target_user_id": userID},
		},
		"matched": true,
	}
	var match []*types.Match
	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &match); err != nil {
		return nil, err
	}
	return match, err
}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("dating").Collection("matches")
}
