package socket

import (
	"context"

	"dating/internal/app/api/types"
	"dating/internal/pkg/glog"
)

type SaveMessage struct {
	message *types.Message `json:"message"`
}
type Repository interface {
	Insert(ctx context.Context, message types.Message) error
}

// chan help save message into db
func saveMessages(sm *chan SaveMessage, r Repository) {
	logger := glog.New().WithField("package", "socket-client")
	for {
		sm, ok := <-*sm
		if !ok {
			logger.Errorf("Error when receiving message to save")
			return
		}
		err := r.Insert(context.Background(), *sm.message)
		if err != nil {
			logger.Errorf("Error when insert message to db", err)

		}
	}
}

// NewSaveMessageChan create a new SaveMessage channel
func NewSaveMessageChan(r Repository) *chan SaveMessage {
	sm := make(chan SaveMessage, 256)
	go saveMessages(&sm, r)
	return &sm
}
