package mailservices

import (
	"context"
	"testing"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	mailpkg "dating/internal/pkg/mail"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryMock struct {
	mock.Mock
}

func (mock *RepositoryMock) FindByEmail(ctx context.Context, email string) (*types.EmailVerification, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.EmailVerification), args.Error(1)
}
func (mock *RepositoryMock) Insert(ctx context.Context, emailVerification types.EmailVerification) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *RepositoryMock) UpdateMailVerified(ctx context.Context, email string) error {
	args := mock.Called()
	return args.Error(0)
}

func TestSendMail(t *testing.T) {
	mockRepo := new(RepositoryMock)

	mockRepo.On("FindByEmail").Return(nil, mongo.ErrNoDocuments).Once()
	mockRepo.On("FindByEmail").Return(&types.EmailVerification{
		ID:          primitive.NewObjectID(),
		Email:       "nguoigiaumat100@gmail.com",
		Code:        "HIUY5d",
		CreatedTime: time.Now(),
		Verified:    false,
	}, nil).Once()
	mockRepo.On("FindByEmail").Return(&types.EmailVerification{
		ID:          primitive.NewObjectID(),
		Email:       "nguoigiaumat100@gmail.com",
		Code:        "HIUY5d",
		CreatedTime: time.Now(),
		Verified:    true,
	}, nil).Once()
	mockRepo.On("Insert").Return(nil)

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	// mock func
	GenCode = func(a int) string {
		return "HIUYas"
	}
	Sendmail = func(content mailpkg.Mail, mails []string, conf *config.Configs) error {
		return nil
	}

	err := testService.SendMail(context.Background(), "nguoigiaumat100@gmail.com")
	assert.Equal(t, nil, err)

	err = testService.SendMail(context.Background(), "nguoigiaumat100@gmail.com")
	assert.Equal(t, nil, err)

	err = testService.SendMail(context.Background(), "nguoigiaumat100@gmail.com")
	assert.NotEqual(t, nil, err)

}

func TestMailVerified(t *testing.T) {
	mockRepo := new(RepositoryMock)

	mockRepo.On("FindByEmail").Return(nil, mongo.ErrNoDocuments).Once()
	mockRepo.On("FindByEmail").Return(&types.EmailVerification{
		ID:          primitive.NewObjectID(),
		Email:       "nguoigiaumat100@gmail.com",
		Code:        "ABCD5D",
		CreatedTime: time.Now(),
		Verified:    true,
	}, nil).Once()
	mockRepo.On("FindByEmail").Return(&types.EmailVerification{
		ID:          primitive.NewObjectID(),
		Email:       "nguoigiaumat100@gmail.com",
		Code:        "ABCD5D",
		CreatedTime: time.Now(),
		Verified:    false,
	}, nil).Twice()
	mockRepo.On("FindByEmail").Return(&types.EmailVerification{
		ID:          primitive.NewObjectID(),
		Email:       "nguoigiaumat100@gmail.com",
		Code:        "ABCD5D",
		CreatedTime: time.Now().Add(-time.Minute * 3),
		Verified:    false,
	}, nil)
	testService := NewService(
		&config.Configs{
			Mail: struct {
				Email    string "mapstructure:\"email\""
				Password string "mapstructure:\"password\""
				Smtp     struct {
					HostMail string "mapstructure:\"host_mail\""
					PortMail string "mapstructure:\"port_mail\""
				} "mapstructure:\"smtp\""
				ConfirmTimeout time.Duration "mapstructure:\"confirm_timeout\""
				SrcTemplate    string        "mapstructure:\"src_template\""
			}{
				ConfirmTimeout: time.Minute * 2,
			},
		},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	err := testService.MailVerified(context.Background(), "nguoigiaumat100@gmail.com", "ABCD5D")
	assert.NotEqual(t, nil, err)

	err = testService.MailVerified(context.Background(), "nguoigiaumat100@gmail.com", "ABCD5D")
	assert.NotEqual(t, nil, err)

	err = testService.MailVerified(context.Background(), "nguoigiaumat100@gmail.com", "ABCD5D")
	assert.Equal(t, nil, err)

	err = testService.MailVerified(context.Background(), "nguoigiaumat100@gmail.com", "ABCD5E")
	assert.NotEqual(t, nil, err)

	err = testService.MailVerified(context.Background(), "nguoigiaumat100@gmail.com", "ABCD5E")
	assert.NotEqual(t, nil, err)
}
