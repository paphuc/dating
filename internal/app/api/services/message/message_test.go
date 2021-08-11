package messageservices

import (
	"context"
	"errors"
	"testing"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/socket"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryMock struct {
	mock.Mock
}

func (mock *RepositoryMock) Insert(ctx context.Context, message types.Message) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *RepositoryMock) FindByIDRoom(ctx context.Context, id string) ([]*types.Message, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*types.Message), args.Error(1)
}

func listMessagesMock() *types.Message {
	return &types.Message{
		ID:          primitive.NewObjectID(),
		RoomID:      primitive.NewObjectID(),
		SenderID:    primitive.NewObjectID(),
		ReceiverID:  primitive.NewObjectID(),
		Content:     "hi",
		Attachments: []string{},
		CreateAt:    time.Now(),
	}
}

func TestGetMessagesByIdRoom(t *testing.T) {
	messages := listMessagesMock()

	mockRepo := new(RepositoryMock)

	mockRepo.On("FindByIDRoom").Return([]*types.Message{messages, messages}, nil).Once()
	mockRepo.On("FindByIDRoom").Return(nil, errors.New("Failed when get list message by id room"))

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	listRooms, err := testService.GetMessagesByIdRoom(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1")
	if err != nil {
		assert.Error(t, err)
	}
	_, err = testService.GetMessagesByIdRoom(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1")
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, 2, len(listRooms))

}
func TestServeWs(t *testing.T) {
	mockRepo := new(RepositoryMock)

	wsServer := socket.NewWebsocketServer()
	go wsServer.Run()
	var conn *websocket.Conn
	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)
	testService.ServeWs(wsServer, conn, "60e3b5d2e1ab4c388ce2d04a", "60e3b5d2e1ab4c388ce2d04a")
	testService.ServeWs(wsServer, conn, "60e3b5d2e1ab4c388ce2d04a", "60e3b5d2e1ab4c388ce2d0422")

}
