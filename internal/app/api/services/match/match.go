package matchservices

import (
	"context"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// Repository is an interface of a user repository
type Repository interface {
	Insert(ctx context.Context, Match types.Match) error
	FindByID(ctx context.Context, id string) (*types.Match, error)
	FindALikeB(ctx context.Context, idUser, idTargetUser string) (*types.Match, error)
	UpdateMatchByID(ctx context.Context, id string) error
	GetListLiked(ctx context.Context, idUser string) ([]*types.Match, error)
	GetListMatched(ctx context.Context, idUser string) ([]*types.Match, error)
	DeleteMatch(ctx context.Context, match types.Match) error
}

// Service is an user service
type Service struct {
	conf   *config.Configs
	em     *config.ErrorMessage
	repo   Repository
	logger glog.Logger
}

// NewService returns a new user service
func NewService(c *config.Configs, e *config.ErrorMessage, r Repository, l glog.Logger) *Service {
	return &Service{
		conf:   c,
		em:     e,
		repo:   r,
		logger: l,
	}
}

// Post basic info user for sign up > A like B
func (s *Service) InsertMatches(ctx context.Context, matchreq types.MatchRequest) (*types.Match, error) {

	// check user B like user A
	matchcheckBA, err := s.repo.FindALikeB(ctx, matchreq.TargetUserID.Hex(), matchreq.UserID.Hex())

	if err != nil {
		// check user A like user B
		matchcheckAB, err := s.repo.FindALikeB(ctx, matchreq.UserID.Hex(), matchreq.TargetUserID.Hex())
		if err != nil {
			match := types.Match{
				ID:           bson.NewObjectId(),
				UserID:       matchreq.UserID,
				TargetUserID: matchreq.TargetUserID,
				Match:        false,
				CreateAt:     time.Now(),
			}
			if err := s.repo.Insert(ctx, match); err != nil {
				s.logger.Errorf("Can't insert user", err)
				return nil, errors.Wrap(err, "Can't insert user")
			}

			matchfind, err := s.repo.FindByID(ctx, match.ID.Hex())
			if err != nil {
				s.logger.Errorf("Can't insert match", err)
				return nil, errors.Wrap(err, "Can't insert match")
			}
			s.logger.Infof("Liked completed", matchreq)
			return matchfind, nil
		}

		// A liked B
		s.logger.Infof("A liked B before", matchreq)
		return matchcheckAB, nil
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

	matchfind, err := s.repo.FindByID(ctx, matchcheckBA.ID.Hex())
	if err != nil {
		s.logger.Errorf("Can't find match", err)
		return nil, errors.Wrap(err, "Can't find match")
	}

	s.logger.Infof("Match completed", matchreq)
	return matchfind, nil

}
