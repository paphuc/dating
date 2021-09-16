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
	FindByIDRoom(ctx context.Context, id string, ps types.PagingNSortingMess) ([]*types.Message, error)
	CountMessage(ctx context.Context, id string) (int64, error)
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
func (s *Service) ServeWs(ctx context.Context, wsServer *socket.WsServer, conn *websocket.Conn, idRoom, idUser string) {

	saveMessagesChan := socket.NewSaveMessageChan(s.repo, s.noti)
	idRoomHex, error := primitive.ObjectIDFromHex(idRoom)

	if error != nil {
		s.logger.Errorc(ctx, "Id room incorrect,it isn't ObjectIdHex %v", error)
		return
	}

	idUserHex, error := primitive.ObjectIDFromHex(idUser)
	if error != nil {
		s.logger.Errorc(ctx, "Id user incorrect,it isn't ObjectIdHex %v", error)
		return
	}

	client := socket.NewClient(conn, wsServer, idRoomHex, idUserHex, saveMessagesChan)

	go client.Write(s.logger)
	go client.Read(s.logger)

	wsServer.Register <- client

}

// method help join client into room message server
func (s *Service) GetMessagesByIdRoom(ctx context.Context, id string, page, size string) (*types.GetListMessageRes, error) {

	var pagingNSorting types.PagingNSortingMess

	if err := pagingNSorting.Init(page, size); err != nil {
		s.logger.Errorc(ctx, "Failed url parameters when get list mess %v", err)
		return nil, err
	}

	count, err := s.repo.CountMessage(ctx, id)
	if err := pagingNSorting.Init(page, size); err != nil {
		s.logger.Errorc(ctx, "Failed when count mess %v", err)
		return nil, err
	}

	listMessages, err := s.repo.FindByIDRoom(ctx, id, pagingNSorting)
	if err != nil {
		s.logger.Errorc(ctx, "Failed when get list message by id room", err)
		return nil, errors.Wrap(err, "Failed when get list message by id room")
	}

	var listMessageRes types.GetListMessageRes

	numberMess := int(count)
	listMessageRes.CurrentPage = pagingNSorting.Page
	listMessageRes.MaxItemsPerPage = int(pagingNSorting.Size)
	listMessageRes.TotalItems = numberMess
	listMessageRes.TotalPages = int(numberMess / pagingNSorting.Size)
	// ex: total: 5, size: 2 => 3 page
	if numberMess%pagingNSorting.Size != 0 {
		listMessageRes.TotalPages += 1
	}

	if pagingNSorting.Size > numberMess {
		listMessageRes.MaxItemsPerPage = numberMess
	}

	listMessageRes.Content = append(listMessageRes.Content, listMessages...)
	if listMessages == nil {
		listMessageRes.Content = []*types.Message{}
	}

	s.logger.Infoc(ctx, "Get list message by id room successfull")

	return &listMessageRes, nil
}
