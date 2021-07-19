package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Photo struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	URLs   []string           `json:"urls" bson:"urls"`
}
