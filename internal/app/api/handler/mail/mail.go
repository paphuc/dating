package mailhandler

import (
	"context"
	"net/http"

	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"
)

type (
	service interface {
		SendMail(ctx context.Context, mail string) error
		MailVerified(ctx context.Context, mail, code string) error
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

// Put handler server message socket HTTP request
func (h *Handler) SendMail(w http.ResponseWriter, r *http.Request) {
	mail := r.URL.Query().Get("email")

	err := h.srv.SendMail(r.Context(), mail)
	if err != nil {
		respond.JSON(w, http.StatusNotFound, h.em.InvalidValue.Request)
		return
	}
	respond.JSON(w, http.StatusOK, h.em.Success)

}

// Put handler server message socket HTTP request
func (h *Handler) MailVerified(w http.ResponseWriter, r *http.Request) {
	mail := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")

	err := h.srv.MailVerified(r.Context(), mail, code)
	if err != nil {
		respond.JSON(w, http.StatusNotFound, h.em.InvalidValue.ValidationFailed)
		return
	}
	respond.JSON(w, http.StatusOK, h.em.Success)

}
