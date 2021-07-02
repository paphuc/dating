package userservices

import (
	"context"

	"dating/internal/app/api/types"
	"dating/internal/pkg/glog"

	"github.com/pkg/errors"
)

// Repository is an interface of a user repository
type Repository interface {
	SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error)
	Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error)
	FindByID(ctx context.Context, id string) (*types.UserResGetInfo, error)
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

// Post basic info user for sign up
func (s *Service) SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error) {
	return s.repo.SignUp(ctx, UserSignUp)
}

// Post basic info user for login
func (s *Service) Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error) {
	return s.repo.Login(ctx, UserLogin)
}

// Get basic info for a user
func (s *Service) FindByID(ctx context.Context, id string) (*types.UserResGetInfo, error) {
	var user *types.UserResGetInfo
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find the given user from database")
	}
	return user, nil
}
