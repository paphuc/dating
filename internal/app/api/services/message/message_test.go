package messageservices

import (
	"context"
	"errors"
	"testing"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	notificationpkg "dating/internal/pkg/notification"
	"dating/internal/pkg/socket"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryMock struct {
	mock.Mock
}
type NotificationServiceMock struct {
	mock.Mock
}

func (mock *NotificationServiceMock) SendNotification(ctx context.Context, id primitive.ObjectID, data notificationpkg.Data, noti notificationpkg.Notification) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *RepositoryMock) Insert(ctx context.Context, message types.Message) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *RepositoryMock) FindByIDRoom(ctx context.Context, id string, ps types.PagingNSortingMess) ([]*types.Message, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*types.Message), args.Error(1)
}
func (mock *RepositoryMock) CountMessage(ctx context.Context, id string) (int64, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(int64), args.Error(1)
}

func listMessagesMock() *types.Message {
	return &types.Message{
		ID:          primitive.NewObjectID(),
		RoomID:      primitive.NewObjectID(),
		Sender:      types.Sender{ID: primitive.NewObjectID()},
		ReceiverID:  primitive.NewObjectID(),
		Content:     "hi",
		Attachments: []string{},
		CreatedAt:   time.Now(),
	}
}

func TestGetMessagesByIdRoom(t *testing.T) {
	messages := listMessagesMock()

	mockRepo := new(RepositoryMock)
	mockNotiService := new(NotificationServiceMock)

	mockNotiService.On("SendNotification").Return(nil)
	mockRepo.On("FindByIDRoom").Return([]*types.Message{messages, messages}, nil).Once()
	mockRepo.On("FindByIDRoom").Return(nil, errors.New("Failed when get list message by id room"))
	mockRepo.On("CountMessage").Return(int64(94), nil).Once()
	mockRepo.On("CountMessage").Return(int64(0), errors.New("Failed when count"))
	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
		mockNotiService,
	)

	listRooms, err := testService.GetMessagesByIdRoom(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", "2", "2")
	if err != nil {
		assert.Error(t, err)
	}
	_, err = testService.GetMessagesByIdRoom(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", "2", "2")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, 2, len(listRooms.Content))

	_, err = testService.GetMessagesByIdRoom(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", "2", "2")
	if err != nil {
		assert.Error(t, err)
	}

}
func TestServeWs(t *testing.T) {
	mockRepo := new(RepositoryMock)
	mockNotiService := new(NotificationServiceMock)

	wsServer := socket.NewWebsocketServer()
	go wsServer.Run()
	var conn *websocket.Conn
	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
		mockNotiService,
	)
	go testService.ServeWs(context.Background(), wsServer, conn, "60e3b5d2e1ab4c388ce2d04a", "60e3b5d2e1ab4c388ce2d04a")
	go testService.ServeWs(context.Background(), wsServer, conn, "60e3b5d2e1ab4c388ce2ds04a", "60e3b5d2e1ab4c388ce2d042")

}
