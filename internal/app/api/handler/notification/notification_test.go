package notificationhandler

import (
	"bytes"
	"context"
	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceMock struct {
	mock.Mock
}

func (mock *ServiceMock) AddDevice(ctx context.Context, noti types.Notification) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *ServiceMock) RemoveDevice(ctx context.Context, noti types.Notification) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *ServiceMock) SendTest(ctx context.Context, id string) error {
	args := mock.Called()
	return args.Error(0)
}

func TestAddDevice(t *testing.T) {
	mockService := new(ServiceMock)
	a, _ := primitive.ObjectIDFromHex("611a1ef8998cb50ada22d162")

	mockService.On("AddDevice").Return(nil).Once()
	mockService.On("AddDevice").Return(errors.New("AddDevice error"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.AddDevice(w, r)
	}))

	defer func() { ts.Close() }()
	body, _ := json.Marshal(types.Notification{
		UserID:      a,
		TokenDevice: "eddndyarQPyVNciA4h_Lj9:APA91bGFEPFwIep3OnUp0zc8DHSesF2QoTauyRsO0YFHPYMqGV1UrbaEX9i-OhgilguQNltvnNLs3iJtrUaaDpz6YwGWAqi5MQQoEN1EgzJU99acSqmGDFrzjyDnKi-kUYRbSi2l03Pa",
	})

	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

func TestRemoveDevice(t *testing.T) {
	mockService := new(ServiceMock)
	a, _ := primitive.ObjectIDFromHex("611a1ef8998cb50ada22d162")

	mockService.On("RemoveDevice").Return(nil).Once()
	mockService.On("RemoveDevice").Return(errors.New("RemoveDevice error"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.RemoveDevice(w, r)
	}))

	defer func() { ts.Close() }()
	body, _ := json.Marshal(types.Notification{
		UserID:      a,
		TokenDevice: "eddndyarQPyVNciA4h_Lj9:APA91bGFEPFwIep3OnUp0zc8DHSesF2QoTauyRsO0YFHPYMqGV1UrbaEX9i-OhgilguQNltvnNLs3iJtrUaaDpz6YwGWAqi5MQQoEN1EgzJU99acSqmGDFrzjyDnKi-kUYRbSi2l03Pa",
	})

	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

func TestTestSend(t *testing.T) {
	mockService := new(ServiceMock)
	mockService.On("SendTest").Return(nil).Once()
	mockService.On("SendTest").Return(errors.New("SendTest error"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.SendTest(w, r)
	}))

	defer func() { ts.Close() }()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	req, err = http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}
