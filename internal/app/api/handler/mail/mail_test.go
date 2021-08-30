package mailhandler

import (
	"context"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

const (
	SendMail     = "SendMail"
	MailVerified = "MailVerified"
)

func (mock *ServiceMock) SendMail(ctx context.Context, mail string) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *ServiceMock) MailVerified(ctx context.Context, mail, code string) error {
	args := mock.Called()
	return args.Error(0)
}

func TestSendMail(t *testing.T) {
	mockService := new(ServiceMock)

	mockService.On(SendMail).Return(nil).Once()
	mockService.On(SendMail).Return(errors.New("Send error"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.SendMail(w, r)
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
	assert.Equal(t, 404, res.StatusCode)
}

func TestMailVerified(t *testing.T) {
	mockService := new(ServiceMock)

	mockService.On(MailVerified).Return(nil).Once()
	mockService.On(MailVerified).Return(errors.New("MailVerified error"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.MailVerified(w, r)
	}))

	defer func() { ts.Close() }()
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
	assert.Equal(t, 404, res.StatusCode)
}
