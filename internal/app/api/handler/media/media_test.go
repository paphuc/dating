package mediahandler

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (mock *ServiceMock) Upload(ctx context.Context, fileBytes []byte) (*types.ImageResponse, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.ImageResponse), args.Error(1)
}
func (mock *ServiceMock) Destroy(ctx context.Context, url string) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *ServiceMock) Asset(ctx context.Context, url string) (*types.ImageResponse, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.ImageResponse), args.Error(1)
}

func TestUpload(t *testing.T) {
	mockService := new(ServiceMock)

	mockService.On("Upload").Return(&types.ImageResponse{Url: "a.jpg"}, nil).Once()
	mockService.On("Upload").Return(nil, errors.New("Cant not upload"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.Upload(w, r)
	}))
	defer func() { ts.Close() }()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "main.go")
	part.Write([]byte(`sample`))
	writer.Close() // <<< important part

	req, err := http.NewRequest("POST", ts.URL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	req, err = http.NewRequest("POST", ts.URL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}
func TestDestroy(t *testing.T) {
	mockService := new(ServiceMock)

	mockService.On("Destroy").Return(errors.New("Cant not destroy"))
	mockService.On("Destroy").Return(nil)

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.Destroy(w, r)
	}))
	defer func() { ts.Close() }()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}
func TestAsset(t *testing.T) {
	mockService := new(ServiceMock)

	mockService.On("Asset").Return(nil, errors.New("Cant not asset")).Once()
	mockService.On("Asset").Return(&types.ImageResponse{Url: "a.jpg"}, nil)

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.Asset(w, r)
	}))
	defer func() { ts.Close() }()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}
