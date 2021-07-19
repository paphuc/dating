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
	"github.com/gorilla/mux"
)

type (
	service interface {
		SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error)
		Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error)
		FindUserById(ctx context.Context, id string) (*types.UserResGetInfo, error)
		UpdateUserByID(ctx context.Context, User types.User) error
		GetListUsers(ctx context.Context, page, size string) (*types.GetListUsersResponse, error)
		GetMatchedUsersByID(ctx context.Context, idUser, matchedParameter string) ([]types.UserResGetInfo, error)
		DisableUserByID(ctx context.Context, idUser, disableParameter string) error
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
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	var userSignup types.UserSignUp

	if err := json.NewDecoder(r.Body).Decode(&userSignup); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(userSignup); err != nil {
		h.logger.Errorf("Failed when validate field userSignup", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
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

	if err := json.NewDecoder(r.Body).Decode(&UserLogin); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if err := validate.Struct(UserLogin); err != nil {
		h.logger.Errorf("Failed when validate field UserLogin", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	user, err := h.srv.Login(r.Context(), UserLogin)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.IncorrectPasswordEmail)
		return
	}

	respond.JSON(w, http.StatusOK, user)
}

// Get handler get information of the user by id
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {

	user, err := h.srv.FindUserById(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, user)
}

// Post handler update information of the user by id
func (h *Handler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {

	var user types.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Errorf("Failed when validate field in method UpdateUserByID", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	if err := validate.Struct(user); err != nil {
		h.logger.Errorf("Failed when validate field in method UpdateUserByID", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.ValidationFailed)
		return
	}

	if error := h.srv.UpdateUserByID(r.Context(), user); error != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, h.em.Success)
}

// Post handler update information of the user by id
func (h *Handler) GetListUsers(w http.ResponseWriter, r *http.Request) {

	pageParameter := r.URL.Query().Get("page")
	sizeParameter := r.URL.Query().Get("size")

	userList, err := h.srv.GetListUsers(r.Context(), pageParameter, sizeParameter)
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, userList)
}

// Get handler get list liked or matched by id idUser HTTP request
func (h *Handler) GetMatchedUsersByID(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["id"]
	matchedParameter := r.URL.Query().Get("matched")

	list, err := h.srv.GetMatchedUsersByID(r.Context(), userID, matchedParameter)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, list)
}
func (h *Handler) DisableUsersByID(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["id"]
	disableParameter := r.URL.Query().Get("disable")

	if err := h.srv.DisableUserByID(r.Context(), userID, disableParameter); err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, h.em.Success)
}
