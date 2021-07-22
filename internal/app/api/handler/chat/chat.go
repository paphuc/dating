package chathandler

import (
	"fmt"
	"net/http"

	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"

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
func New(c *config.Configs, e *config.ErrorMessage, s service, l glog.Logger) *Handler {
	return &Handler{
		conf:   c,
		em:     e,
		srv:    s,
		logger: l,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Put handler post insert match HTTP request
func (h *Handler) WS(w http.ResponseWriter, r *http.Request) {

	// fmt.Printf("w's type is %T\n", w)
	matchedParameter := r.URL.Query().Get("room")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf(err.Error())
	}

	fmt.Println("New Client joined the hub!")
	fmt.Println(ws)
	fmt.Println(matchedParameter)
	respond.JSON(w, http.StatusOK, nil)
}
