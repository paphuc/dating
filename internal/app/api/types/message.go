package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sender struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}
type Message struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	RoomID      primitive.ObjectID `json:"room_id" bson:"room_id"`
	Sender      Sender             `json:"sender" bson:"sender"`
	ReceiverID  primitive.ObjectID `json:"receiver_id" bson:"receiver_id"`
	Content     string             `json:"content" bson:"content"`
	Attachments []string           `json:"attachments" bson:"attachments"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
