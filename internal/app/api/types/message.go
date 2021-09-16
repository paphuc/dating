package types

import (
	"strconv"
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
type PaginationMessage struct {
	TotalItems      int `json:"totalItems"`
	TotalPages      int `json:"totalPages"`
	CurrentPage     int `json:"currentPage"`
	MaxItemsPerPage int `json:"maxItemsPerPage"`
}
type ListMessageRes struct {
	Content []*Message `json:"content"`
}
type GetListMessageRes struct {
	PaginationMessage
	ListMessageRes
}

type PagingNSortingMess struct {
	Size int `json:"size" default:"100"`
	Page int `json:"page" default:"1"`
}

func (ps *PagingNSortingMess) Init(page, size string) error {

	if size == "" && page == "" {
		ps.Size = 100
		ps.Page = 1
		return nil
	}

	if size == "" {
		size = "100" // default size = 1
	}

	sizeInt, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return err
	}

	if page == "" {
		page = "1" // default page = 1
	}

	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return err
	}

	ps.Page = int(pageInt)
	ps.Size = int(sizeInt)

	return nil
}
