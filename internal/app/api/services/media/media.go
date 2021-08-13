package mediaservices

import (
	"context"
	"dating/internal/app/api/types"
	"dating/internal/app/config"
	cloudinarypkg "dating/internal/pkg/cloudinary"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/uuid"
)

// Repository is an interface of a match repository
type Repository interface {
}
type Cloudinary interface {
	UploadFile(ctx context.Context, fileBytes []byte, name string) (string, error)
	DestroyFile(ctx context.Context, name string) (string, error)
	AssetFile(ctx context.Context, name string) (string, error)
}

// Service is an match service
type Service struct {
	conf   *config.Configs
	em     *config.ErrorMessage
	repo   Repository
	logger glog.Logger
	cloud  Cloudinary
}

// NewService returns a new match service
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

// Post basic, user match someone
func (s *Service) Upload(ctx context.Context, fileBytes []byte) (*types.ImageResponse, error) {

	src, err := s.cloud.UploadFile(ctx, fileBytes, uuid.New())
	if err != nil {
		return nil, err
	}

	return &types.ImageResponse{Url: src}, nil
}

// Post basic, user match someone
func (s *Service) Destroy(ctx context.Context, url string) error {
	_, error := s.cloud.DestroyFile(ctx, url)
	return error
}

// Post basic, user match someone
func (s *Service) Asset(ctx context.Context, url string) (*types.ImageResponse, error) {
	src, err := s.cloud.AssetFile(ctx, url)
	if err != nil {
		return nil, err
	}
	return &types.ImageResponse{Url: src}, nil
}
