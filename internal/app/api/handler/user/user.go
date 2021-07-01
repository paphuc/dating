package userhandler

import (
	"context"
	"encoding/json"
	"net/http"

	"dating/internal/app/api/types"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"
)

type (
	service interface {
		SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error)
		Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error)
	}
	// Handler is user web handler
	Handler struct {
		srv    service
		logger glog.Logger
	}
)

// New returns new res api user handler
func New(s service, l glog.Logger) *Handler {
	return &Handler{
		srv:    s,
		logger: l,
	}
}

// Post handler  post sign up HTTP request
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	var userSignup types.UserSignUp
	err := json.NewDecoder(r.Body).Decode(&userSignup)

	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	user, err := h.srv.SignUp(r.Context(), userSignup)
	if err != nil {
		respond.Error(w, err, http.StatusConflict)
		return
	}

	respond.JSON(w, http.StatusOK, user)
}

// Post handler  login HTTP request
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var UserLogin types.UserLogin
	err := json.NewDecoder(r.Body).Decode(&UserLogin)

	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	user, err := h.srv.Login(r.Context(), UserLogin)
	if err != nil {
		respond.Error(w, err, http.StatusUnauthorized)
		return
	}
	respond.JSON(w, http.StatusOK, user)
}
