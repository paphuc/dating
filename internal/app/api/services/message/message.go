package messageservices

import (
	"context"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	notification "dating/internal/pkg/notification"
	socket "dating/internal/pkg/socket"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository is an interface of a message repository
type Repository interface {
	Insert(ctx context.Context, message types.Message) error
	FindByIDRoom(ctx context.Context, id string) ([]*types.Message, error)
}
type NotificationService interface {
	SendNotification(ctx context.Context, id primitive.ObjectID, data notification.Data, noti notification.Notification) error
}

// Service is an message service
type Service struct {
	conf   *config.Configs
	em     *config.ErrorMessage
	repo   Repository
	logger glog.Logger
	noti   NotificationService
}

// NewService returns a new message service
func NewService(c *config.Configs, e *config.ErrorMessage, r Repository, l glog.Logger, n NotificationService) *Service {
	return &Service{
		conf:   c,
		em:     e,
		repo:   r,
		logger: l,
		noti:   n,
	}
}

// method help join client into room message server
func (s *Service) ServeWs(wsServer *socket.WsServer, conn *websocket.Conn, idRoom, idUser string) {

	saveMessagesChan := socket.NewSaveMessageChan(s.repo, s.noti)
	idRoomHex, error := primitive.ObjectIDFromHex(idRoom)

	if error != nil {
		s.logger.Errorf("Id room incorrect,it isn't ObjectIdHex %v", error)
		return
	}

	idUserHex, error := primitive.ObjectIDFromHex(idUser)
	if error != nil {
		s.logger.Errorf("Id user incorrect,it isn't ObjectIdHex %v", error)
		return
	}

	client := socket.NewClient(conn, wsServer, idRoomHex, idUserHex, saveMessagesChan)

	go client.Write(s.logger)
	go client.Read(s.logger)

	wsServer.Register <- client

}

// method help join client into room message server
func (s *Service) GetMessagesByIdRoom(ctx context.Context, id string) ([]*types.Message, error) {

	listMessages, err := s.repo.FindByIDRoom(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed when get list message by id room", err)
		return nil, errors.Wrap(err, "Failed when get list message by id room")
	}
	s.logger.Infof("Get list message by id room successfull")

	if listMessages == nil {
		return []*types.Message{}, nil
	}

	return listMessages, nil
}
