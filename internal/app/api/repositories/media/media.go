package media

import (
	"context"

	"dating/internal/app/api/types"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
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

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("dating").Collection("media")
}
