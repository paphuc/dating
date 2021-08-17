package notification

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

// This method helps , update notification
func (r *MongoRepository) Insert(ctx context.Context, noti types.Notification) error {

	filter := bson.M{
		"user_id":      noti.UserID,
		"token_device": noti.TokenDevice,
	}
	updated := bson.M{
		"$set": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection().UpdateOne(ctx, filter, updated, opts)
	return err

}

// This method helps insert notification
func (r *MongoRepository) Delete(ctx context.Context, noti types.Notification) error {

	filter := bson.M{
		"user_id":      noti.UserID,
		"token_device": noti.TokenDevice,
	}

	_, err := r.collection().DeleteOne(ctx, filter)
	return err

}

func (r *MongoRepository) Find(ctx context.Context, id primitive.ObjectID) ([]*types.Notification, error) {

	filter := bson.M{
		"user_id": id,
	}

	var result []*types.Notification
	opts := options.Find()

	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, err

}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("dating").Collection("notification")
}
