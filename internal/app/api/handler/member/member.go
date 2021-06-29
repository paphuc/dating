package memberhandler

import (
	"context"
	"net/http"

	"dating/internal/app/types"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"

	"github.com/gorilla/mux"
)

type (
	service interface {
		Get(ctx context.Context, id string) (*types.Member, error)
	}

	// Handler is member web handler
	Handler struct {
		srv    service
		logger glog.Logger
	}
)

// New return new rest api member handler
func New(s service, l glog.Logger) *Handler {
	return &Handler{
		srv:    s,
		logger: l,
	}
}

// Get handle get member HTTP request
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	member, err := h.srv.Get(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, member)
}
