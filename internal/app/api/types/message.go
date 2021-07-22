package types

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Message struct {
	ID          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	RoomID      bson.ObjectId `json:"room_id" bson:"_room_id"`
	SenderID    bson.ObjectId `json:"sender_id" bson:"sender_id"`
	Content     string        `json:"content" bson:"content"`
	Attachments string        `json:"attachments" bson:"attachments"`
	CreateAt    time.Time     `json:"created_at" bson:"created_at"`
}
type Rooms struct {
	ID           bson.ObjectId   `json:"_id" bson:"_id,omitempty"`
	Participants []bson.ObjectId `json:"participants" bson:"participants"`
}
