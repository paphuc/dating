package notificationservices

import (
	"context"
	"encoding/json"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	notification "dating/internal/pkg/notification"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository is an interface of a notification repository
type Repository interface {
	Insert(ctx context.Context, noti types.Notification) error
	Delete(ctx context.Context, noti types.Notification) error
	Find(ctx context.Context, id primitive.ObjectID) ([]*types.Notification, error)
}

// Service is an notification service
type Service struct {
	conf   *config.Configs
	em     *config.ErrorMessage
	repo   Repository
	logger glog.Logger
}

// NewService returns a new notification service
func NewService(c *config.Configs, e *config.ErrorMessage, r Repository, l glog.Logger) *Service {
	return &Service{
		conf:   c,
		em:     e,
		repo:   r,
		logger: l,
	}
}

var (
	PushNotification = notification.PushNotification
)

// method help add Device in notification
func (s *Service) AddDevice(ctx context.Context, noti types.Notification) error {

	err := s.repo.Insert(context.Background(), noti)
	if err != nil {
		s.logger.Errorc(ctx, "Can't AddDevice: %v", err)
		return err
	}
	s.logger.Errorc(ctx, "AddDevice completed %v", noti)
	return nil
}

// method help Remove Device in notification
func (s *Service) RemoveDevice(ctx context.Context, noti types.Notification) error {
	err := s.repo.Delete(context.Background(), noti)
	if err != nil {
		s.logger.Errorc(ctx, "Can't RemoveDevice: %v", err)
		return err
	}
	s.logger.Errorc(ctx, "RemoveDevice completed %v", noti)
	return nil
}

// method help test send notification
func (s *Service) SendTest(ctx context.Context, id string) error {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	list, err := s.repo.Find(context.Background(), ID)
	if err != nil {
		s.logger.Errorc(ctx, "Find notification failed: %v", err)
		return err
	}
	for _, noti := range list {
		if time.Now().Sub(noti.CreateAt) > s.conf.Jwt.Duration {
			err := s.repo.Delete(context.Background(), *noti)
			if err != nil {
				s.logger.Errorc(ctx, "Can't remove notification: %v", err)
				return err
			}
			s.logger.Errorc(ctx, "Remove token devices completed: %v", noti)
		} else {
			payLoad := notification.NotificationPayLoad{
				RegistrationIds: []string{noti.TokenDevice},
				Data: notification.Data{
					Content: "Test push notification",
				},
				Foreground: true,
				Notification: notification.Notification{
					Title: "notification",
					Body:  "test notification body",
				},
			}
			plByte, _ := json.Marshal(payLoad)
			result := make(chan error, 1)

			PushNotification(s.conf, plByte, result)

			value, ok := <-result
			if ok {
				if value != nil {
					s.logger.Errorc(ctx, "Can't send notification: %v", value)
				}
			}
			defer close(result)
		}
	}
	s.logger.Infoc(ctx, "send notification completed")
	return nil
}

// method help send notification
func (s *Service) SendNotification(ctx context.Context, id primitive.ObjectID, dataPayLoad notification.Data, notiPayLoad notification.Notification) error {
	list, err := s.repo.Find(context.Background(), id)
	if err != nil {
		s.logger.Errorc(ctx, "Find notification failed: %v", err)
		return err
	}
	for _, noti := range list {
		if time.Now().Sub(noti.CreateAt) > s.conf.Jwt.Duration {
			err := s.repo.Delete(context.Background(), *noti)
			if err != nil {
				s.logger.Errorc(ctx, "Can't remove notification: %v", err)
				return err
			}
			s.logger.Errorc(ctx, "Remove token devices completed: %v", noti)
		} else {
			payLoad := notification.NotificationPayLoad{
				RegistrationIds: []string{noti.TokenDevice},
				Data:            dataPayLoad,
				Notification:    notiPayLoad,
				Foreground:      true,
			}
			plByte, _ := json.Marshal(payLoad)
			result := make(chan error, 1)

			PushNotification(s.conf, plByte, result)

			value, ok := <-result
			if ok {
				if value != nil {
					s.logger.Errorc(ctx, "Can't send notification: %v", value)
				}
			}
			defer close(result)
		}
	}
	s.logger.Infoc(ctx, "send notification completed")
	return nil
}
