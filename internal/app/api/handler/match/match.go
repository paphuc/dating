package matchhandler

import (
	"context"
	"encoding/json"
	"net/http"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type (
	service interface {
		InsertMatch(ctx context.Context, Match types.MatchRequest) (*types.Match, error)
		DeleteMatch(ctx context.Context, matchreq types.MatchRequest) error
		FindRoomsByUserId(ctx context.Context, id string) ([]types.MatchRoomResponse, error)
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

// Put handler post insert match HTTP request
func (h *Handler) InsertMatch(w http.ResponseWriter, r *http.Request) {

	var matchRequest types.MatchRequest

	if err := json.NewDecoder(r.Body).Decode(&matchRequest); err != nil {
		h.logger.Errorf("Failed when NewDecoder matchRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(matchRequest); err != nil {
		h.logger.Errorf("Failed when validate field matchRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	match, err := h.srv.InsertMatch(r.Context(), matchRequest)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, match)
}

// Del handler unMatch or unlike by matched HTTP request
func (h *Handler) DeleteMatched(w http.ResponseWriter, r *http.Request) {

	var unmatchRequest types.MatchRequest

	if err := json.NewDecoder(r.Body).Decode(&unmatchRequest); err != nil {
		h.logger.Errorf("Failed when NewDecoder matchRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(unmatchRequest); err != nil {
		h.logger.Errorf("Failed when validate field matchRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	err := h.srv.DeleteMatch(r.Context(), unmatchRequest)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, h.em.Success)
}

// Put handler get list room by user id
func (h *Handler) GetRoomsByUserId(w http.ResponseWriter, r *http.Request) {

	roomList, err := h.srv.FindRoomsByUserId(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, roomList)
}
