package media

import (
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

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("dating").Collection("media")
}
