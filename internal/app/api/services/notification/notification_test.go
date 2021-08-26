package notificationservices

import (
	"context"
	"testing"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/notification"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryMock struct {
	mock.Mock
}

func (mock RepositoryMock) Insert(ctx context.Context, noti types.Notification) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock RepositoryMock) Delete(ctx context.Context, noti types.Notification) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock RepositoryMock) Find(ctx context.Context, id primitive.ObjectID) ([]*types.Notification, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*types.Notification), args.Error(1)
}

func NotiMock() *types.Notification {
	id, _ := primitive.ObjectIDFromHex("611a1ef8998cb50ada22d162")
	return &types.Notification{
		UserID:      id,
		TokenDevice: "eddndyarQPyVNciA4h_Lj9:APA91bGFEPFwIep3OnUp0zc8DHSesF2QoTauyRsO0YFHPYMqGV1UrbaEX9i-OhgilguQNltvnNLs3iJtrUaaDpz6YwGWAqi5MQQoEN1EgzJU99acSqmGDFrzjyDnKi-kUYRbSi2l03Pa",
		CreateAt:    time.Now(),
	}
}

func TestInsert(t *testing.T) {

	PushNotification = func(conf *config.Configs, payLoad []byte, result chan<- error) {
		result <- nil
	}
	mockRepo := new(RepositoryMock)
	mockRepo.On("Insert").Return(nil).Once()
	mockRepo.On("Insert").Return(errors.New("Insert failed"))
	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	err := testService.AddDevice(context.Background(), *NotiMock())
	err1 := testService.AddDevice(context.Background(), *NotiMock())
	mockRepo.AssertExpectations(t)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, err1)
}

func TestRemove(t *testing.T) {

	PushNotification = func(conf *config.Configs, payLoad []byte, result chan<- error) {
		result <- nil
	}
	mockRepo := new(RepositoryMock)
	mockRepo.On("Delete").Return(nil).Once()
	mockRepo.On("Delete").Return(errors.New("Insert failed"))
	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	err := testService.RemoveDevice(context.Background(), *NotiMock())
	err1 := testService.RemoveDevice(context.Background(), *NotiMock())
	mockRepo.AssertExpectations(t)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, err1)
}

func TestSendNotification(t *testing.T) {
	a, _ := primitive.ObjectIDFromHex("611a1ef8998cb50ada22d162")
	PushNotification = func(conf *config.Configs, payLoad []byte, result chan<- error) {
		result <- nil
	}
	mockRepo := new(RepositoryMock)
	noti := NotiMock()
	noti1 := NotiMock()
	noti1.CreateAt = time.Now().AddDate(-int(22+1), 0, 0)

	mockRepo.On("Find").Return([]*types.Notification{noti, noti, noti, noti}, nil).Once()
	mockRepo.On("Find").Return([]*types.Notification{noti, noti1, noti, noti}, nil).Once()
	mockRepo.On("Find").Return(nil, errors.New("find failed"))

	mockRepo.On("Delete").Return(nil)

	testService := NewService(
		&config.Configs{
			Jwt: struct {
				Duration time.Duration "mapstructure:\"duration\""
			}{
				Duration: time.Hour * 29,
			},
		},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)
	err := testService.SendNotification(context.Background(), a, notification.Data{}, notification.Notification{})
	err1 := testService.SendNotification(context.Background(), a, notification.Data{}, notification.Notification{})
	err2 := testService.SendNotification(context.Background(), a, notification.Data{}, notification.Notification{})

	mockRepo.AssertExpectations(t)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, err1)
	assert.NotEqual(t, nil, err2)
}

func TestTestSend(t *testing.T) {
	PushNotification = func(conf *config.Configs, payLoad []byte, result chan<- error) {
		result <- nil
	}
	mockRepo := new(RepositoryMock)
	noti := NotiMock()
	noti1 := NotiMock()
	noti1.CreateAt = time.Now().AddDate(-int(22+1), 0, 0)

	mockRepo.On("Find").Return([]*types.Notification{noti, noti, noti, noti}, nil).Once()
	mockRepo.On("Find").Return([]*types.Notification{noti, noti1, noti, noti}, nil).Once()
	mockRepo.On("Find").Return(nil, errors.New("find failed"))

	mockRepo.On("Delete").Return(nil)

	testService := NewService(
		&config.Configs{
			Jwt: struct {
				Duration time.Duration "mapstructure:\"duration\""
			}{
				Duration: time.Hour * 29,
			},
		},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)
	err := testService.TestSend(context.Background(), "611a1ef8998cb50ada22d162")
	err1 := testService.TestSend(context.Background(), "611a1ef8998cb50ada22d162")
	err2 := testService.TestSend(context.Background(), "611a1ef8998cb50ada22d162")

	mockRepo.AssertExpectations(t)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, err1)
	assert.NotEqual(t, nil, err2)
}
