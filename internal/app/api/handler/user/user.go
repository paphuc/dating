package userhandler

import (
	"context"
	"encoding/json"
	"net/http"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type (
	service interface {
		SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error)
		Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error)
		FindByID(ctx context.Context, id string) (*types.UserResGetInfo, error)
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
		respond.JSON(w, http.StatusInternalServerError, config.EM.Invalid_value.Request)
		return
	}

	user, err := h.srv.SignUp(r.Context(), userSignup)
	if err != nil {
		respond.JSON(w, http.StatusConflict, config.EM.Invalid_value.Email_exists)
		return
	}

	respond.JSON(w, http.StatusOK, user)
}

// Post handler  login HTTP request
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var UserLogin types.UserLogin
	err := json.NewDecoder(r.Body).Decode(&UserLogin)

	if err != nil {
		respond.JSON(w, http.StatusUnauthorized, config.EM.Invalid_value.Incorrect_password_email)
		return
	}

	user, err := h.srv.Login(r.Context(), UserLogin)
	if err != nil {
		respond.JSON(w, http.StatusUnauthorized, config.EM.Invalid_value.Incorrect_password_email)
		return
	}
	respond.JSON(w, http.StatusOK, user)
}

// Post handler get your own information
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {

	token, ok := r.Context().Value("props").(map[string]interface{})
	if !ok {
		respond.Error(w, errors.New("failed to"), http.StatusUnauthorized)
	}

	idUser := token["_id"].(string)
	user, err := h.srv.FindByID(r.Context(), idUser)
	if err != nil {
		respond.JSON(w, http.StatusUnauthorized, config.EM.Invalid_value.Request)
	}

	respond.JSON(w, http.StatusOK, user)
}

// Post handler get infomation of the user by id
func (h *Handler) FindById(w http.ResponseWriter, r *http.Request) {

	user, err := h.srv.FindByID(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.Error(w, err, http.StatusUnauthorized)
	}

	respond.JSON(w, http.StatusOK, user)
}
