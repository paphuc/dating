package matchservices

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryMock struct {
	mock.Mock
}

func (mock RepositoryMock) Insert(ctx context.Context, Match types.Match) error {
	return nil
}

func (mock RepositoryMock) FindByID(ctx context.Context, id string) (*types.Match, error) {
	return nil, nil
}

func (mock RepositoryMock) FindALikeB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.Match), args.Error(1)
}

func (mock RepositoryMock) UpdateMatchByID(ctx context.Context, id string) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock RepositoryMock) GetListLiked(ctx context.Context, idUser string) ([]*types.Match, error) {
	return nil, nil
}

func (mock RepositoryMock) GetListMatched(ctx context.Context, idUser string) ([]*types.Match, error) {
	return nil, nil
}

func (mock RepositoryMock) DeleteMatch(ctx context.Context, id string) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock RepositoryMock) UpsertMatch(ctx context.Context, match types.Match) error {
	args := mock.Called()
	return args.Error(0)

}

func (mock RepositoryMock) CheckAB(ctx context.Context, idUser, idTargetUser string, matched bool) (*types.Match, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.Match), args.Error(1)
}

func (mock RepositoryMock) FindAMatchB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.Match), args.Error(1)
}

func (mock RepositoryMock) FindRoomsByUserId(ctx context.Context, idUser string) ([]*types.MatchRoomResponse, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*types.MatchRoomResponse), args.Error(1)
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
	a, err := primitive.ObjectIDFromHex("60e3b5d2e1ab4c388ce2d04a")
	if err != nil {
		assert.Error(t, err)
		return
	}
	b, err := primitive.ObjectIDFromHex("60e3b5dbe1ab4c388ce2d04c")
	if err != nil {
		assert.Error(t, err)
		return
	}

	mockRepo := new(RepositoryMock)

	mockRepo.On("FindALikeB").Return(&types.Match{
		ID:           primitive.NewObjectID(),
		UserID:       b,
		TargetUserID: a,
		Matched:      false,
		CreateAt:     time.Now(),
	}, nil).Once()
	mockRepo.On("UpdateMatchByID").Return(nil).Once()
	mockRepo.On("FindALikeB").Return(nil, errors.New("err")).Once()
	mockRepo.On("UpsertMatch").Return(nil).Once()

	mockRepo.On("FindALikeB").Return(nil, errors.New("err")).Once()
	mockRepo.On("UpsertMatch").Return(errors.New("err")).Once()

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	result, err := testService.InsertMatch(context.Background(), types.MatchRequest{
		Matched:      false,
		UserID:       a,
		TargetUserID: b,
	})
	if err != nil {
		assert.Error(t, err)
		return
	}
	result2, err := testService.InsertMatch(context.Background(), types.MatchRequest{
		Matched:      false,
		UserID:       a,
		TargetUserID: b,
	})
	if err != nil {
		assert.Error(t, err)
		return
	}
	result3, err := testService.InsertMatch(context.Background(), types.MatchRequest{
		Matched:      false,
		UserID:       a,
		TargetUserID: b,
	})
	if err != nil {
		assert.Error(t, err)
		return
	}

	mockRepo.AssertExpectations(t)
	assert.Equal(t, b, result.UserID)
	assert.Equal(t, true, result2.Matched)
	assert.Equal(t, nil, result3)

}

func TestDeleteMatch(t *testing.T) {
	match, err := matchMock()
	if err != nil {
		assert.Error(t, err)
	}

	mockRepo := new(RepositoryMock)

	mockRepo.On("FindAMatchB").Return(match, nil).Times(1)
	mockRepo.On("CheckAB").Return(match, nil).Times(1)
	mockRepo.On("DeleteMatch").Return(nil).Times(2)

	mockRepo.On("FindAMatchB").Return(match, nil)
	mockRepo.On("CheckAB").Return(nil, errors.New("err"))
	mockRepo.On("DeleteMatch").Return(errors.New("err"))

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	err = testService.DeleteMatch(context.Background(), types.MatchRequest{
		Matched:      true,
		UserID:       match.UserID,
		TargetUserID: match.TargetUserID,
	})
	err1 := testService.DeleteMatch(context.Background(), types.MatchRequest{
		Matched:      false,
		UserID:       match.UserID,
		TargetUserID: match.TargetUserID,
	})

	err2 := testService.DeleteMatch(context.Background(), types.MatchRequest{
		Matched:      true,
		UserID:       match.UserID,
		TargetUserID: match.TargetUserID,
	})
	err12 := testService.DeleteMatch(context.Background(), types.MatchRequest{
		Matched:      false,
		UserID:       match.UserID,
		TargetUserID: match.TargetUserID,
	})
	fmt.Println(err, err1, err2, err12)

	mockRepo.AssertExpectations(t)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, err1)
	assert.NotEqual(t, nil, err2)
	assert.NotEqual(t, nil, err12)
}

func TestFindRoomsByUserId(t *testing.T) {
	mockRepo := new(RepositoryMock)

	usersID1 := primitive.NewObjectID()
	usersID2 := primitive.NewObjectID()

	rooms := &types.MatchRoomResponse{
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

	mockRepo.On("FindRoomsByUserId").Return([]*types.MatchRoomResponse{rooms, rooms, rooms}, nil)

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	listRooms, err := testService.FindRoomsByUserId(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1")
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.Equal(t, 3, len(listRooms))
}
