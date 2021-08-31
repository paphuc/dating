package mediaservices

import (
	"context"
	"dating/internal/app/api/types"
	"dating/internal/app/config"
	cloudinarypkg "dating/internal/pkg/cloudinary"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/uuid"
)

// Repository is an interface of a media repository
type Repository interface {
}

// Repository is an interface of a Cloudinary
type Cloudinary interface {
	UploadFile(ctx context.Context, fileBytes []byte, name string) (string, error)
	DestroyFile(ctx context.Context, name string) (string, error)
	AssetFile(ctx context.Context, name string) (string, error)
}

// Service is an media service
type Service struct {
	conf   *config.Configs
	em     *config.ErrorMessage
	repo   Repository
	logger glog.Logger
	cloud  Cloudinary
}

// NewService returns a new media service
func NewService(c *config.Configs, e *config.ErrorMessage, r Repository, l glog.Logger) *Service {
	cld, _ := cloudinarypkg.New(c.Cloudinary.URL)
	return &Service{
		conf:   c,
		em:     e,
		repo:   r,
		logger: l,
		cloud:  cld,
	}
}

// Post upload media
func (s *Service) Upload(ctx context.Context, fileBytes []byte) (*types.MediaResponse, error) {

	src, err := s.cloud.UploadFile(ctx, fileBytes, uuid.New())
	if err != nil {
		s.logger.Errorc(ctx, "Failed when upload %v", err)
		return nil, err
	}
	s.logger.Infoc(ctx, "Upload successfully %v", src)
	return &types.MediaResponse{Url: src}, nil
}

// Post del media
func (s *Service) Destroy(ctx context.Context, url string) error {
	_, error := s.cloud.DestroyFile(ctx, url)
	s.logger.Infoc(ctx, "Destroy with err: %v", error)
	return error
}

// Post get media
func (s *Service) Asset(ctx context.Context, url string) (*types.MediaResponse, error) {
	src, err := s.cloud.AssetFile(ctx, url)
	if err != nil {
		s.logger.Errorc(ctx, "Asset when upload %v", err)
		return nil, err
	}
	s.logger.Infoc(ctx, "Asset successfully %v", src)
	return &types.MediaResponse{Url: src}, nil
}
