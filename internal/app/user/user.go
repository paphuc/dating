package user

import (
	"context"
	"dating/internal/app/types"
	"dating/internal/pkg/glog"
)

// Repository is an interface of a user repository
type Repository interface {
	SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error)
	Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error)
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
func (s *Service) Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error) {
	return s.repo.Login(ctx, UserLogin)
}
