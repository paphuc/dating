package mediahandler

import (
	"context"
	"io/ioutil"
	"net/http"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"

	"github.com/gorilla/mux"
)

type (
	service interface {
		Upload(ctx context.Context, fileBytes []byte) (*types.MediaResponse, error)
		Destroy(ctx context.Context, url string) error
		Asset(ctx context.Context, url string) (*types.MediaResponse, error)
	}
	// Handler is media web handler
	Handler struct {
		conf   *config.Configs
		em     *config.ErrorMessage
		srv    service
		logger glog.Logger
	}
)

// New returns new res api media handler
func New(c *config.Configs, e *config.ErrorMessage, s service, l glog.Logger) *Handler {
	return &Handler{
		conf:   c,
		em:     e,
		srv:    s,
		logger: l,
	}
}

// Put handler update media
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := r.FormFile("file")
	if err != nil {
		h.logger.Errorf("file failed: %v", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}
	defer file.Close()
	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		h.logger.Errorf("ReadAll file: %v", err)
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	media, err := h.srv.Upload(r.Context(), fileBytes)
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}

	respond.JSON(w, http.StatusOK, media)
}

// Del handler delete media
func (h *Handler) Destroy(w http.ResponseWriter, r *http.Request) {
	err := h.srv.Destroy(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}
	respond.JSON(w, http.StatusOK, h.em.Success)
}

// Get handler get media
func (h *Handler) Asset(w http.ResponseWriter, r *http.Request) {
	media, err := h.srv.Asset(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		respond.JSON(w, http.StatusBadRequest, h.em.InvalidValue.Request)
		return
	}
	respond.JSON(w, http.StatusOK, media)
}
