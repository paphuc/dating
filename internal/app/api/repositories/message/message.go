package message

import (
	"context"

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

// This method helps insert message
func (r *MongoRepository) Insert(ctx context.Context, message types.Message) error {
	_, err := r.collection().InsertOne(context.TODO(), message)
	return err
}

// This method helps get message by id room
func (r *MongoRepository) FindByIDRoom(ctx context.Context, id string) ([]*types.Message, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"room_id": objectID,
	}

	var result []*types.Message
	opts := options.Find()

	cursor, err := r.collection().Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, err
}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("dating").Collection("message_test")
}
