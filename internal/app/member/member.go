package member

import (
	"context"
	"dating/internal/app/types"
	"dating/internal/pkg/glog"
)

// Repository is an interface of a member repository
type Repository interface {
	FindByID(ctx context.Context, id string) (*types.Member, error)
}

// Service is an member service
type Service struct {
	repo   Repository
	logger glog.Logger
}

// NewService return a new member service
func NewService(r Repository, l glog.Logger) *Service {
	return &Service{
		repo:   r,
		logger: l,
	}
}

// Get return given member by his/her id
func (s *Service) Get(ctx context.Context, id string) (*types.Member, error) {
	return s.repo.FindByID(ctx, id)
}
