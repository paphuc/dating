package notificationservices

import (
	"context"
	"encoding/json"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	notificationpkg "dating/internal/pkg/notification"

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
	PushNotification = notificationpkg.PushNotification
)

// method help add Device in notification
func (s *Service) AddDevice(ctx context.Context, noti types.Notification) error {

	err := s.repo.Insert(context.Background(), noti)
	if err != nil {
		s.logger.Errorf("Can't AddDevice: %v", err)
		return err
	}
	s.logger.Errorf("AddDevice completed %v", noti)
	return nil
}

// method help Remove Device in notification
func (s *Service) RemoveDevice(ctx context.Context, noti types.Notification) error {
	err := s.repo.Delete(context.Background(), noti)
	if err != nil {
		s.logger.Errorf("Can't RemoveDevice: %v", err)
		return err
	}
	s.logger.Errorf("RemoveDevice completed %v", noti)
	return nil
}

// method help test send notification
func (s *Service) TestSend(ctx context.Context, id string) error {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	list, err := s.repo.Find(context.Background(), ID)
	if err != nil {
		s.logger.Errorf("Find notification failed: %v", err)
		return err
	}
	listN := s.convertPointerArrayToArrayNotification(list)
	for _, notification := range listN {
		if time.Now().Sub(notification.CreateAt) > s.conf.Jwt.Duration {
			err := s.repo.Delete(context.Background(), notification)
			if err != nil {
				s.logger.Errorf("Can't remove notification: %v", err)
				return err
			}
			s.logger.Errorf("Remove token devices completed: %v", notification)
		} else {
			payLoad := notificationpkg.NotificationPayLoad{
				RegistrationIds: []string{notification.TokenDevice},
				Data: notificationpkg.Data{
					Content: "Test push notification",
				},
				Foreground: true,
				Notification: notificationpkg.Notification{
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
					s.logger.Errorf("Can't send notification: %v", value)
				}
			}
			defer close(result)
		}
	}
	s.logger.Infof("send notification completed")
	return nil
}

// method help send notification
func (s *Service) SendNotification(ctx context.Context, id primitive.ObjectID, data notificationpkg.Data, noti notificationpkg.Notification) error {
	list, err := s.repo.Find(context.Background(), id)
	if err != nil {
		s.logger.Errorf("Find notification failed: %v", err)
		return err
	}
	listN := s.convertPointerArrayToArrayNotification(list)
	for _, notification := range listN {
		if time.Now().Sub(notification.CreateAt) > s.conf.Jwt.Duration {
			err := s.repo.Delete(context.Background(), notification)
			if err != nil {
				s.logger.Errorf("Can't remove notification: %v", err)
				return err
			}
			s.logger.Errorf("Remove token devices completed: %v", notification)
		} else {
			payLoad := notificationpkg.NotificationPayLoad{
				RegistrationIds: []string{notification.TokenDevice},
				Data:            data,
				Notification:    noti,
				Foreground:      true,
			}
			plByte, _ := json.Marshal(payLoad)
			result := make(chan error, 1)

			PushNotification(s.conf, plByte, result)

			value, ok := <-result
			if ok {
				if value != nil {
					s.logger.Errorf("Can't send notification: %v", value)
				}
			}
			defer close(result)
		}
	}
	s.logger.Infof("send notification completed")
	return nil
}

// convert []*types.Notification to []types.Notification - if empty return []
func (s *Service) convertPointerArrayToArrayNotification(list []*types.Notification) []types.Notification {
	listN := []types.Notification{}
	for _, mgs := range list {
		listN = append(listN, *mgs)
	}
	return listN
}
