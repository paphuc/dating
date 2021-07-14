package matchservices

import (
	"context"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"

	"github.com/pkg/errors"
)

// Repository is an interface of a match repository
type Repository interface {
	Insert(ctx context.Context, Match types.Match) error
	FindByID(ctx context.Context, id string) (*types.Match, error)
	FindALikeB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error)
	UpdateMatchByID(ctx context.Context, id string) error
	GetListLiked(ctx context.Context, idUser string) ([]*types.Match, error)
	GetListMatched(ctx context.Context, idUser string) ([]*types.Match, error)
	DeleteMatch(ctx context.Context, id string) error
	UpsertMatch(ctx context.Context, match types.Match) error
	CheckAB(ctx context.Context, idUser, idTargetUser string, matched bool) (*types.Match, error)
	FindAMatchB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error)
}

// Service is an match service
type Service struct {
	conf   *config.Configs
	em     *config.ErrorMessage
	repo   Repository
	logger glog.Logger
}

// NewService returns a new match service
func NewService(c *config.Configs, e *config.ErrorMessage, r Repository, l glog.Logger) *Service {
	return &Service{
		conf:   c,
		em:     e,
		repo:   r,
		logger: l,
	}
}

// Post basic, user match someone
func (s *Service) InsertMatch(ctx context.Context, matchreq types.MatchRequest) (*types.Match, error) {

	// check user B like user A
	matchcheckBA, err := s.repo.FindALikeB(ctx, matchreq.TargetUserID.Hex(), matchreq.UserID.Hex())

	if err != nil {
		match := types.Match{
			UserID:       matchreq.UserID,
			TargetUserID: matchreq.TargetUserID,
			Match:        false,
			CreateAt:     time.Now(),
		}

		err := s.repo.UpsertMatch(ctx, match)
		if err != nil {
			s.logger.Errorf("Can't update match", err)
			return nil, errors.Wrap(err, "Can't update match")
		}

		s.logger.Infof("A liked B before", matchreq)
		return &match, nil
	}
	// B liked A
	if matchcheckBA.Match {
		s.logger.Infof("B, A matched before", matchreq)
		return matchcheckBA, nil
	}
	if err := s.repo.UpdateMatchByID(ctx, matchcheckBA.ID.Hex()); err != nil {
		s.logger.Errorf("Can't update match", err)
		return nil, errors.Wrap(err, "Can't update match")
	}

	matchcheckBA.Match = true

	s.logger.Infof("Match completed", matchreq)
	return matchcheckBA, nil
}

// Post basic help user unlike someone
func (s *Service) Unlike(ctx context.Context, matchreq types.MatchRequest) error {
	// check user A like user B
	matchcheckAB, err := s.repo.CheckAB(ctx, matchreq.UserID.Hex(), matchreq.TargetUserID.Hex(), false)
	if err != nil {
		s.logger.Errorf("UserA haven't liked B", err)
		return err
	}
	if err := s.repo.DeleteMatch(ctx, matchcheckAB.ID.Hex()); err != nil {
		s.logger.Errorf("Can't del like", err)
		return err
	}

	s.logger.Infof("Unlike completed", matchreq)
	return nil
}

// Post basic help user unMatch someone
func (s *Service) Unmatched(ctx context.Context, matchreq types.MatchRequest) error {
	// check user A matched user B
	matchcheckAB, err := s.repo.FindAMatchB(ctx, matchreq.UserID.Hex(), matchreq.TargetUserID.Hex())
	if err != nil {
		s.logger.Errorf("A B have not matched before", err)
		return err
	}
	if err := s.repo.DeleteMatch(ctx, matchcheckAB.ID.Hex()); err != nil {
		s.logger.Errorf("Can't del match", err)
		return err
	}

	s.logger.Infof("Unmatched completed", matchreq)
	return nil
}
