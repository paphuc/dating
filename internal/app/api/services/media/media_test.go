package mediaservices

import (
	"context"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

type CloudinaryMock struct {
	mock.Mock
}

func (mock *CloudinaryMock) UploadFile(ctx context.Context, fileBytes []byte, name string) (string, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return "", args.Error(1)
	}
	return result.(string), args.Error(1)
}

func (mock *CloudinaryMock) DestroyFile(ctx context.Context, name string) (string, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return "", args.Error(1)
	}
	return result.(string), args.Error(1)
}

func (mock *CloudinaryMock) AssetFile(ctx context.Context, name string) (string, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return "", args.Error(1)
	}
	return result.(string), args.Error(1)
}

func TestService_Upload(t *testing.T) {
	mockRepo := new(RepositoryMock)
	mockCloud := new(CloudinaryMock)

	mockCloud.On("UploadFile").Return("", errors.New("Cant not upload")).Once()
	mockCloud.On("UploadFile").Return("a.jpg", nil)

	testService := Service{
		conf:   &config.Configs{},
		em:     &config.ErrorMessage{},
		repo:   mockRepo,
		logger: glog.New(),
		cloud:  mockCloud,
	}

	_, err := testService.Upload(context.Background(), []byte{})

	assert.Error(t, err)

	media, err := testService.Upload(context.Background(), []byte{})
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.Equal(t, nil, err)
	assert.Equal(t, "a.jpg", media.Url)
}
func TestService_Destroy(t *testing.T) {
	mockRepo := new(RepositoryMock)
	mockCloud := new(CloudinaryMock)

	mockCloud.On("DestroyFile").Return("", errors.New("Cant not destroy"))
	mockCloud.On("DestroyFile").Return("ok", nil)

	testService := Service{
		conf:   &config.Configs{},
		em:     &config.ErrorMessage{},
		repo:   mockRepo,
		logger: glog.New(),
		cloud:  mockCloud,
	}

	err := testService.Destroy(context.Background(), "aaaa")
	assert.Error(t, err)

	err = testService.Destroy(context.Background(), "aaaa")
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.Equal(t, nil, err)
}
func TestService_Asset(t *testing.T) {
	mockRepo := new(RepositoryMock)
	mockCloud := new(CloudinaryMock)

	mockCloud.On("AssetFile").Return("", errors.New("Cant not upload")).Once()
	mockCloud.On("AssetFile").Return("a.jpg", nil)

	testService := Service{
		conf:   &config.Configs{},
		em:     &config.ErrorMessage{},
		repo:   mockRepo,
		logger: glog.New(),
		cloud:  mockCloud,
	}

	_, err := testService.Asset(context.Background(), "a")

	assert.Error(t, err)

	media, err := testService.Asset(context.Background(), "a")
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.Equal(t, nil, err)
	assert.Equal(t, "a.jpg", media.Url)
}
