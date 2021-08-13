package mailservices

import (
	"context"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	mailpkg "dating/internal/pkg/mail"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository is an interface of a mail repository
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*types.EmailVerification, error)
	Insert(ctx context.Context, emailVerification types.EmailVerification) error
}

// Service is an mail service
type Service struct {
	conf   *config.Configs
	em     *config.ErrorMessage
	repo   Repository
	logger glog.Logger
}

// NewService returns a new mail service
func NewService(c *config.Configs, e *config.ErrorMessage, r Repository, l glog.Logger) *Service {
	return &Service{
		conf:   c,
		em:     e,
		repo:   r,
		logger: l,
	}
}

var (
	GenCode  = mailpkg.GenCode
	Sendmail = mailpkg.Sendmail
)

// method help verify mail users
func (s *Service) MailVerified(ctx context.Context, mail, code string) error {

	emailExists, err := s.repo.FindByEmail(ctx, mail)
	if err != nil {
		s.logger.Errorf("Email not found %v", err)
		return errors.Wrap(err, "Email not found")
	}

	if emailExists.Verified {
		s.logger.Errorf("MailVerified has true")
		return errors.Wrap(errors.New("MailVerified failed"), "Email has confirmed")
	}

	if time.Now().Sub(emailExists.CreatedTime) > s.conf.Mail.ConfirmTimeout {
		s.logger.Errorf("Code expired")
		return errors.Wrap(errors.New("MailVerified failed"), "Code expired")
	}
	if emailExists.Code != code {
		s.logger.Errorf("Code incorrect")
		return errors.Wrap(errors.New("MailVerified failed"), "Code incorrect")
	}
	s.logger.Infof("MailVerified has completed")
	return nil

}

// method help find and send  mail user
func (s *Service) SendMail(ctx context.Context, email string) error {

	emailExists, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err := s.GenNumberSendAndUpsert(ctx, email); err != nil {
				s.logger.Errorf("Error when gen code %v", err)
				return err
			}
			return nil
		}
	}

	if emailExists.Verified {
		s.logger.Errorf("Send mail fail. Email had confirm")
		return errors.Wrap(errors.New("Send mail fail"), "Email had confirm")
	}

	if err := s.GenNumberSendAndUpsert(ctx, email); err != nil {
		s.logger.Errorf("Error when gen code %v", err)
		return err
	}
	s.logger.Infof("Send email completed")
	return nil
}

// method help gen code and send code to mail user
func (s *Service) GenNumberSendAndUpsert(ctx context.Context, mail string) error {

	code := GenCode(6)

	//upsert db
	if err := s.repo.Insert(ctx, types.EmailVerification{
		Email:    mail,
		Code:     code,
		Verified: false,
	}); err != nil {
		return err
	}

	//send email
	if err := Sendmail(mailpkg.Mail{
		Subject: "Mail Confirm",
		Body:    code,
	}, []string{mail}, s.conf); err != nil {
		return errors.Wrap(err, "Send mail fail")
	}

	s.logger.Infof("Send email completed %v", mail)
	return nil
}
