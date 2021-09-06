package messagehandler

import (
	"context"
	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	socket "dating/internal/pkg/socket"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceMock struct {
	mock.Mock
}

func (mock *ServiceMock) ServeWs(ctx context.Context, wsServer *socket.WsServer, conn *websocket.Conn, idRoom, idUser string) {

}
func (mock *ServiceMock) GetMessagesByIdRoom(ctx context.Context, id string, page, size string) (*types.GetListMessageRes, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.GetListMessageRes), args.Error(1)
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
	mockService := new(ServiceMock)
	message := listMessagesMock()
	mockService.On("GetMessagesByIdRoom").Return(&types.GetListMessageRes{
		types.PaginationMessage{
			TotalItems:      95,
			TotalPages:      24,
			CurrentPage:     2,
			MaxItemsPerPage: 4,
		},
		types.ListMessageRes{
			Content: []*types.Message{message, message, message},
		},
	}, nil).Once()
	mockService.On("GetMessagesByIdRoom").Return(&types.GetListMessageRes{}, errors.New("Can't like or match"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.GetMessagesByIdRoom(w, r)
	}))

	defer func() { ts.Close() }()
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body_res, _ := io.ReadAll(res.Body)
	var body_mock *types.GetListMessageRes
	json.Unmarshal([]byte(body_res), &body_mock)
	assert.Equal(t, 3, len(body_mock.Content))

	req, err = http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, res.StatusCode)
}

func TestServeWs(t *testing.T) {
	mockService := new(ServiceMock)
	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.ServeWs(w, r)
	}))
	defer func() { ts.Close() }()

	u := "ws" + strings.TrimPrefix(ts.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer ws.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}
