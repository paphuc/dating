package pictureservices

import (
	"context"

	"dating/internal/app/api/types"
	"dating/internal/pkg/glog"
)

// Repository is an interface of a user repository
type Repository interface {
	Insert(ctx context.Context, User types.Photo) error
}

// Service is an user service
type Service struct {
	repo   Repository
	logger glog.Logger
}

// NewService returns a new user service
func NewService(r Repository, l glog.Logger) *Service {
	return &Service{
		repo:   r,
		logger: l,
	}
}
