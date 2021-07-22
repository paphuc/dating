package chathandler

import (
	"fmt"
	"net/http"

	"dating/internal/app/config"
	"dating/internal/pkg/glog"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type (
	service interface {
	}
	// Handler is match web handler
	Handler struct {
		conf   *config.Configs
		em     *config.ErrorMessage
		srv    service
		logger glog.Logger
	}
)

var (
	validate = validator.New()
)

// New returns new res api match handler
func New(c *config.Configs, e *config.ErrorMessage, l glog.Logger) *Handler {
	return &Handler{
		conf: c,
		em:   e,
		// srv:    s,
		logger: l,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Put handler post insert match HTTP request
func (h *Handler) WS(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorc(r.Context(), err.Error())
	}
	fmt.Println(ws)

	defer ws.Close()
}
