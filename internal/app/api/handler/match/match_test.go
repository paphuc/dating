package matchhandler

import (
	"bytes"
	"context"
	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceMock struct {
	mock.Mock
}

func (mock *ServiceMock) InsertMatch(ctx context.Context, Match types.MatchRequest) (*types.Match, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.Match), args.Error(1)
}
func (mock *ServiceMock) DeleteMatch(ctx context.Context, matchreq types.MatchRequest) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *ServiceMock) FindRoomsByUserId(ctx context.Context, id string) ([]types.MatchRoomResponse, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]types.MatchRoomResponse), args.Error(1)
}

func matchMock() (*types.Match, error) {
	a, err := primitive.ObjectIDFromHex("60e3b5d2e1ab4c388ce2d04a")
	if err != nil {
		return nil, err
	}
	b, err := primitive.ObjectIDFromHex("60e3b5dbe1ab4c388ce2d04c")
	if err != nil {
		return nil, err
	}

	return &types.Match{
		ID:           primitive.NewObjectID(),
		UserID:       a,
		TargetUserID: b,
		Matched:      false,
		CreateAt:     time.Now(),
	}, nil
}
func TestInsertMatch(t *testing.T) {

	match, _ := matchMock()

	mockService := new(ServiceMock)
	mockService.On("InsertMatch").Return(match, nil).Once()
	mockService.On("InsertMatch").Return(nil, errors.New("Can't like or match"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.InsertMatch(w, r)
	}))
	defer func() { ts.Close() }()

	a, _ := primitive.ObjectIDFromHex("60e3b5d2e1ab4c388ce2d04a")
	b, _ := primitive.ObjectIDFromHex("60e3b5dbe1ab4c388ce2d04c")
	matchRequest := types.MatchRequest{
		UserID:       a,
		TargetUserID: b,
	}

	body, _ := json.Marshal(``)
	req, err := http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	body, _ = json.Marshal(matchRequest)
	req, err = http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	body, _ = json.Marshal(matchRequest)
	req, err = http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

}
func TestDeleteMatched(t *testing.T) {

	mockService := new(ServiceMock)
	mockService.On("DeleteMatch").Return(nil).Once()
	mockService.On("DeleteMatch").Return(errors.New("Can't like or match"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.DeleteMatched(w, r)
	}))
	defer func() { ts.Close() }()

	a, _ := primitive.ObjectIDFromHex("60e3b5d2e1ab4c388ce2d04a")
	b, _ := primitive.ObjectIDFromHex("60e3b5dbe1ab4c388ce2d04c")
	matchRequest := types.MatchRequest{
		UserID:       a,
		TargetUserID: b,
		Matched:      true,
	}

	body, _ := json.Marshal(``)
	req, err := http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	body, _ = json.Marshal(matchRequest)
	req, err = http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	body, _ = json.Marshal(matchRequest)
	req, err = http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

}

func TestGetRoomsByUserId(t *testing.T) {

	usersID1 := primitive.NewObjectID()
	usersID2 := primitive.NewObjectID()

	rooms := types.MatchRoomResponse{
		ID: primitive.NewObjectID(),
		User: []types.UserResGetInfoInRoom{
			{
				ID:     usersID1,
				Name:   "Phuc",
				Avatar: "src",
				Gender: "Male",
			},
			{
				ID:     usersID2,
				Name:   "Huynh",
				Avatar: "src",
				Gender: "Female",
			},
		},
		LastMessage: &types.Message{
			ID:          primitive.NewObjectID(),
			RoomID:      primitive.NewObjectID(),
			SenderID:    usersID1,
			ReceiverID:  usersID2,
			Content:     "Hi",
			Attachments: []string{},
			CreateAt:    time.Now(),
		},
	}

	mockService := new(ServiceMock)
	mockService.On("FindRoomsByUserId").Return([]types.MatchRoomResponse{rooms, rooms, rooms}, nil).Once()
	mockService.On("FindRoomsByUserId").Return([]types.MatchRoomResponse{}, errors.New("Can't like or match"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.GetRoomsByUserId(w, r)
	}))
	defer func() { ts.Close() }()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body_res, _ := io.ReadAll(res.Body)
	var body_mock []types.MatchRoomResponse
	json.Unmarshal([]byte(body_res), &body_mock)
	assert.Equal(t, 3, len(body_mock))

	req, err = http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, res.StatusCode)

}
