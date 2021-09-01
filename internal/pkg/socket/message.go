package socket

import (
	"context"

	"dating/internal/app/api/types"
	"dating/internal/pkg/glog"
	notificationpkg "dating/internal/pkg/notification"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SaveMessage struct {
	message *types.Message `json:"message"`
}
type Repository interface {
	Insert(ctx context.Context, message types.Message) error
}
type NotificationService interface {
	SendNotification(ctx context.Context, id primitive.ObjectID, data notificationpkg.Data, noti notificationpkg.Notification) error
}

// chan help save message into db
func saveMessages(sm *chan SaveMessage, r Repository, noti NotificationService) {
	logger := glog.New().WithField("package", "socket-client")
	for {
		sm, ok := <-*sm
		if !ok {
			logger.Errorc(context.Background(), "Error when receiving message to save")
			return
		}
		err := r.Insert(context.Background(), *sm.message)
		if err != nil {
			logger.Errorc(context.Background(), "Error when insert message to db %v", err)
		}
		noti.SendNotification(context.Background(), sm.message.ReceiverID, notificationpkg.Data{}, notificationpkg.Notification{
			Body:  sm.message.Content,
			Title: "new message",
		})
	}
}

// NewSaveMessageChan create a new SaveMessage channel
func NewSaveMessageChan(r Repository, noti NotificationService) *chan SaveMessage {
	sm := make(chan SaveMessage, 256)
	go saveMessages(&sm, r, noti)
	return &sm
}
