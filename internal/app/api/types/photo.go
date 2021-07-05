package types

import "github.com/globalsign/mgo/bson"

type Photo struct {
	ID     bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	UserID bson.ObjectId `json:"user_id" bson:"user_id"`
	URLs   []string      `json:"urls" bson:"urls"`
}
