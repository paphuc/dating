package userhandler

import (
	"context"
	"encoding/json"
	"net/http"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

type (
	service interface {
		SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error)
		Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error)
		FindUserById(ctx context.Context, id string) (*types.UserResGetInfo, error)
		UpdateUserByID(ctx context.Context, User types.User) error
	}
	// Handler is user web handler
	Handler struct {
		conf   *config.Configs
		em     *config.ErrorMessage
		srv    service
		logger glog.Logger
	}
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
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	var userSignup types.UserSignUp
	err := json.NewDecoder(r.Body).Decode(&userSignup)

	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	user, err := h.srv.SignUp(r.Context(), userSignup)

	if err != nil {
		respond.JSON(w, http.StatusConflict, h.em.InvalidValue.EmailExists)
		return
	}

	respond.JSON(w, http.StatusOK, user)
}

// Post handler  login HTTP request
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var UserLogin types.UserLogin
	err := json.NewDecoder(r.Body).Decode(&UserLogin)

	if err != nil {
		respond.JSON(w, http.StatusUnauthorized, h.em.InvalidValue.IncorrectPasswordEmail)
		return
	}

	user, err := h.srv.Login(r.Context(), UserLogin)
	if err != nil {
		respond.JSON(w, http.StatusUnauthorized, h.em.InvalidValue.IncorrectPasswordEmail)
		return
	}
	respond.JSON(w, http.StatusOK, user)
}

// Get handler get your own information
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {

	token, ok := r.Context().Value("props").(map[string]interface{})
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, h.em.InvalidValue.FailedAuthentication)
	}
	// get id from token
	idUser := token["_id"].(string)

	user, err := h.srv.FindUserById(r.Context(), idUser)
	if err != nil {
		respond.JSON(w, http.StatusUnauthorized, h.em.InvalidValue.Request)
	}

	respond.JSON(w, http.StatusOK, user)
}

// Get handler get information of the user by id
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {

	user, err := h.srv.FindUserById(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.JSON(w, http.StatusUnauthorized, h.em.InvalidValue.Request)
	}

	respond.JSON(w, http.StatusOK, user)
}

// Post handler update information of the user by id
func (h *Handler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {

	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	token, ok := r.Context().Value("props").(map[string]interface{})
	if !ok {
		respond.JSON(w, http.StatusUnauthorized, h.em.InvalidValue.FailedAuthentication)
	}

	// get id,email from token
	idUser := token["_id"].(string)
	user.ID = bson.ObjectIdHex(idUser)
	user.Email = token["email"].(string)

	error := h.srv.UpdateUserByID(r.Context(), user)

	if error != nil {
		respond.JSON(w, http.StatusUnauthorized, h.em.InvalidValue.Request)
	}
	respond.JSON(w, http.StatusOK, h.em.Success)
}
