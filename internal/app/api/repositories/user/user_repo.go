package user

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

// this method helps insert user
func (r *MongoRepository) Insert(ctx context.Context, user types.User) error {
	_, err := r.collection().InsertOne(ctx, user)
	return err
}

// this method helps get user with email
func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	var user *types.User
	err := r.collection().FindOne(ctx, bson.M{"email": email, "disable": false}).Decode(&user)
	return user, err
}

// This method helps get basic info user
func (r *MongoRepository) FindByID(ctx context.Context, id string) (*types.UserResGetInfo, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user *types.UserResGetInfo
	err = r.collection().FindOne(ctx, bson.M{"_id": objectID, "disable": false}).Decode(&user)

	return user, err
}

//  This method helps update info user
func (r *MongoRepository) UpdateUserByID(ctx context.Context, user types.User) error {
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

	_, err := r.collection().UpdateByID(ctx, user.ID, updatedUser)
	return err
}

// This method helps Enable/Disable account
func (r *MongoRepository) DisableUserByID(ctx context.Context, idUser string, disable bool) error {
	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return err
	}
	disableUpdate := bson.M{"$set": bson.M{
		"disable": disable,
	}}
	_, err = r.collection().UpdateByID(ctx, userID, disableUpdate)
	return err
}

// This method helps get all user by page
func (r *MongoRepository) GetListUsers(ctx context.Context, ps types.PagingNSorting) ([]*types.UserResGetInfo, error) {
	query := bson.M{
		"disable": false,
		"birthday": bson.M{
			"$gte": ps.Filter.AgeRange.Gte,
			"$lt":  ps.Filter.AgeRange.Lt,
		},
		"gender": bson.M{
			"$in": ps.Filter.Gender,
		},
	}
	var result []*types.UserResGetInfo
	opts := options.Find()
	opts.SetSkip(int64((ps.Page - 1) * ps.Size))
	opts.SetLimit(int64(ps.Size))
	cursor, err := r.collection().Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, err
}

// This method helps count number users in collection
func (r *MongoRepository) CountUser(ctx context.Context, ps types.PagingNSorting) (int64, error) {
	query := bson.M{
		"disable": false,
		"birthday": bson.M{
			"$gte": ps.Filter.AgeRange.Gte,
			"$lt":  ps.Filter.AgeRange.Lt,
		},
		"gender": bson.M{
			"$in": ps.Filter.Gender,
		},
	}
	return r.collection().CountDocuments(ctx, query)
}

// this method help get list matched include info
func (r *MongoRepository) GetListMatchedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error) {
	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return nil, err
	}
	query := []bson.M{
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"user_id": userID},
				bson.M{"target_user_id": userID},
			},
			"matched": true,
		}},
		{"$project": bson.M{
			"targer_id": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$user_id", userID}},
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
	cursor, err := r.client.Database("dating").Collection("matches").Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &listMatched); err != nil {
		return nil, err
	}
	return listMatched, err
}

// this method help get list liked include info
func (r *MongoRepository) GetListlikedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error) {
	userID, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"user_id": userID,
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
	cursor, err := r.client.Database("dating").Collection("matches").Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &listMatched); err != nil {
		return nil, err
	}
	return listMatched, err
}

// This method helps find people who had matched, liked
func (r *MongoRepository) IgnoreIdUsers(ctx context.Context, id string) ([]primitive.ObjectID, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := []bson.M{
		{"$match": bson.M{
			"$or": []interface{}{
				bson.M{"user_id": userID},
				bson.M{
					"matched":        true,
					"target_user_id": userID,
				},
			},
		}},
		{"$project": bson.M{
			"targer_id": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$user_id", userID}},
					"$target_user_id", "$user_id"},
			},
		},
		},
	}

	cursor, err := r.client.Database("dating").Collection("matches").Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}

	var result []*struct {
		ID        primitive.ObjectID `json:"_id" bson:"_id"`
		Targer_id primitive.ObjectID `json:"targer_id" bson:"targer_id"`
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	ignoreIds := []primitive.ObjectID{}

	for _, c := range result {
		ignoreIds = append(ignoreIds, c.Targer_id)
	}
	ignoreIds = append(ignoreIds, userID) // yourself

	return ignoreIds, nil
}

// This method helps get all user by page ignore self and people who had matched, liked
func (r *MongoRepository) GetListUsersAvailable(ctx context.Context, ignoreIds []primitive.ObjectID, ps types.PagingNSorting) ([]*types.UserResGetInfo, error) {

	query := bson.M{
		"disable": false,
		"birthday": bson.M{
			"$gte": ps.Filter.AgeRange.Gte,
			"$lt":  ps.Filter.AgeRange.Lt,
		},
		"gender": bson.M{
			"$in": ps.Filter.Gender,
		},
		"_id": bson.M{
			"$nin": ignoreIds,
		},
	}
	var result []*types.UserResGetInfo
	opts := options.Find()
	opts.SetSkip(int64((ps.Page - 1) * ps.Size))
	opts.SetLimit(int64(ps.Size))
	cursor, err := r.collection().Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, err
}

// This method helps count number users in collection, ignore self and people who had matched, liked
func (r *MongoRepository) CountUserUsersAvailable(ctx context.Context, ignoreIds []primitive.ObjectID, ps types.PagingNSorting) (int64, error) {
	query := bson.M{
		"disable": false,
		"birthday": bson.M{
			"$gte": ps.Filter.AgeRange.Gte,
			"$lt":  ps.Filter.AgeRange.Lt,
		},
		"gender": bson.M{
			"$in": ps.Filter.Gender,
		},
		"_id": bson.M{
			"$nin": ignoreIds,
		},
	}
	return r.collection().CountDocuments(ctx, query)
}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("dating").Collection("users")
}
