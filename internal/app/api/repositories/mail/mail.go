package mail

import (
	"context"
	"time"

	"dating/internal/app/api/types"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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
func (r *MongoRepository) Insert(ctx context.Context, emailVerification types.EmailVerification) error {

	filter := bson.M{
		"email": emailVerification.Email,
	}
	updated := bson.M{
		"$set": bson.M{
			"code":         emailVerification.Code,
			"verified":     emailVerification.Verified,
			"created_time": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection().UpdateOne(ctx, filter, updated, opts)
	return err

}

// this method helps get user with email
func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (*types.EmailVerification, error) {
	var result *types.EmailVerification
	err := r.collection().FindOne(ctx, bson.M{"email": email}).Decode(&result)
	return result, err
}

// this method help get update verify email in mail collection
func (r *MongoRepository) UpdateMailVerified(ctx context.Context, email string) error {
	filter := bson.M{
		"email": email,
	}
	updated := bson.M{
		"$set": bson.M{
			"verified": true,
		},
	}
	_, err := r.collection().UpdateOne(ctx, filter, updated)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) collection() *mongo.Collection {
	return r.client.Database("dating").Collection("mails")
}
