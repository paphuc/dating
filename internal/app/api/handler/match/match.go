package userhandler

import (
	"context"
	"encoding/json"
	"net/http"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"

	"github.com/go-playground/validator/v10"
)

type (
	service interface {
		InsertMatches(ctx context.Context, Match types.MatchRequest) (*types.Match, error)
		UnMatched(ctx context.Context, Match types.MatchRequest) error
	}
	// Handler is user web handler
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

// New returns new res api user handler
func New(c *config.Configs, e *config.ErrorMessage, s service, l glog.Logger) *Handler {
	return &Handler{
		conf:   c,
		em:     e,
		srv:    s,
		logger: l,
	}
}

// Post handler  post sign up HTTP request
func (h *Handler) InsertMatches(w http.ResponseWriter, r *http.Request) {

	var matchRequest types.MatchRequest

	if err := json.NewDecoder(r.Body).Decode(&matchRequest); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(matchRequest); err != nil {
		h.logger.Errorf("Failed when validate field matchRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	match, err := h.srv.InsertMatches(r.Context(), matchRequest)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, match)
}

// Post handler  post sign up HTTP request
func (h *Handler) UnMatch(w http.ResponseWriter, r *http.Request) {

	var matchRequest types.MatchRequest

	if err := json.NewDecoder(r.Body).Decode(&matchRequest); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(matchRequest); err != nil {
		h.logger.Errorf("Failed when validate field matchRequest", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	err := h.srv.UnMatched(r.Context(), matchRequest)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, h.em.Success)
}
