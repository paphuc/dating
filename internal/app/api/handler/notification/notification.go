package notificationhandler

import (
	"context"
	"encoding/json"
	"net/http"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"

	"github.com/gorilla/mux"
)

type (
	service interface {
		AddDevice(ctx context.Context, noti types.Notification) error
		RemoveDevice(ctx context.Context, noti types.Notification) error
		TestSend(ctx context.Context, id string) error
	}
	// Handler is message web handler
	Handler struct {
		conf   *config.Configs
		em     *config.ErrorMessage
		srv    service
		logger glog.Logger
	}
)

// New returns new res api message handler
func New(c *config.Configs, e *config.ErrorMessage, s service, l glog.Logger) *Handler {
	return &Handler{
		conf:   c,
		em:     e,
		srv:    s,
		logger: l,
	}
}
func (h *Handler) AddDevice(w http.ResponseWriter, r *http.Request) {

	var noti types.Notification

	if err := json.NewDecoder(r.Body).Decode(&noti); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}
	err := h.srv.AddDevice(r.Context(), noti)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}
	respond.JSON(w, http.StatusOK, h.em.Success)
}

// Put handler get list message of room
func (h *Handler) RemoveDevice(w http.ResponseWriter, r *http.Request) {

	var noti types.Notification

	if err := json.NewDecoder(r.Body).Decode(&noti); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}
	err := h.srv.RemoveDevice(r.Context(), noti)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}
	respond.JSON(w, http.StatusOK, h.em.Success)
}

// Put handler get list message of room
func (h *Handler) TestSend(w http.ResponseWriter, r *http.Request) {

	err := h.srv.TestSend(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}
	respond.JSON(w, http.StatusOK, h.em.Success)
}
