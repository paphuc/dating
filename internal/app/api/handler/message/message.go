package messagehandler

import (
	"context"
	"net/http"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"
	socket "dating/internal/pkg/socket"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type (
	service interface {
		ServeWs(ctx context.Context, wsServer *socket.WsServer, conn *websocket.Conn, idRoom, idUser string)
		GetMessagesByIdRoom(ctx context.Context, id string) ([]*types.Message, error)
	}
	// Handler is message web handler
	Handler struct {
		conf   *config.Configs
		em     *config.ErrorMessage
		srv    service
		logger glog.Logger
	}
)

var (
	wsServer = socket.NewWebsocketServer()

	socketBufferSize  = 2048
	messageBufferSize = 256

	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// New returns new res api message handler
func New(c *config.Configs, e *config.ErrorMessage, s service, l glog.Logger) *Handler {

	go wsServer.Run()

	return &Handler{
		conf:   c,
		em:     e,
		srv:    s,
		logger: l,
	}
}

// Put handler server message socket HTTP request
func (h *Handler) ServeWs(w http.ResponseWriter, r *http.Request) {
	idRoom := r.URL.Query().Get("room")
	idUser := r.URL.Query().Get("user")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		h.logger.Errorc(r.Context(), "Can't create ServeWs for client %v", err.Error())
		return
	}

	h.srv.ServeWs(r.Context(), wsServer, conn, idRoom, idUser)

	h.logger.Infoc(r.Context(), "New Client joined the room! %v", idRoom)
}

// Put handler get list message of room
func (h *Handler) GetMessagesByIdRoom(w http.ResponseWriter, r *http.Request) {

	messagesList, err := h.srv.GetMessagesByIdRoom(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, messagesList)
}
