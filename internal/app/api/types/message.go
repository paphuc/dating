package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	RoomID      primitive.ObjectID `json:"room_id" bson:"room_id"`
	SenderID    primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	ReceiverID  primitive.ObjectID `json:"receiver_id" bson:"receiver_id"`
	Content     string             `json:"content" bson:"content"`
	Attachments []string           `json:"attachments" bson:"attachments"`
	CreateAt    time.Time          `json:"created_at" bson:"created_at"`
}
