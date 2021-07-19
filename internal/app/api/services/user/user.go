package userservices

import (
	"context"
	"strconv"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/jwt"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// Repository is an interface of a user repository
type Repository interface {
	FindByID(ctx context.Context, id string) (*types.UserResGetInfo, error)
	FindByEmail(ctx context.Context, email string) (*types.User, error)
	Insert(ctx context.Context, User types.User) error
	UpdateUserByID(ctx context.Context, User types.User) error
	GetListUsers(ctx context.Context, ps types.PagingNSorting) ([]*types.UserResGetInfo, error)
	CountUser(ctx context.Context) (int, error)
	GetListlikedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error)
	GetListMatchedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error)
	DisableUserByID(ctx context.Context, idUser string, disable bool) error
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

// Post basic info user for sign up
func (s *Service) SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error) {

	if _, err := s.repo.FindByEmail(ctx, UserSignUp.Email); err == nil {
		s.logger.Errorf("Email email exits", err)
		return nil, errors.Wrap(errors.New("Email email exits"), "Email exits, can't insert user")
	}

	UserSignUp.Password, _ = jwt.HashPassword(UserSignUp.Password)

	user := types.User{
		ID:       bson.NewObjectId(),
		Name:     UserSignUp.Name,
		Email:    UserSignUp.Email,
		Password: UserSignUp.Password,
		Disable:  false,
		CreateAt: time.Now(),
		UpdateAt: time.Now()}

	if err := s.repo.Insert(ctx, user); err != nil {
		s.logger.Errorf("Can't insert user", err)
		return nil, errors.Wrap(err, "Can't insert user")
	}

	var tokenString string
	tokenString, err := jwt.GenToken(types.UserFieldInToken{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, s.conf.Jwt.Duration)

	if err != nil {
		s.logger.Errorf("Can't gen token after insert", err)
		return nil, errors.Wrap(err, "Can't insert user")
	}
	s.logger.Infof("Register completed", UserSignUp)

	return &types.UserResponseSignUp{
		Name:  UserSignUp.Name,
		Email: UserSignUp.Email,
		Token: tokenString}, nil

}

// Post basic info user for login
func (s *Service) Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error) {

	user, err := s.repo.FindByEmail(ctx, UserLogin.Email)
	if err != nil {
		s.logger.Errorf("Not found email exits", err)
		return nil, errors.Wrap(errors.New("Not found email exits"), "Email not exists, can't find user")
	}

	if !jwt.IsCorrectPassword(UserLogin.Password, user.Password) {
		s.logger.Errorf("Password incorrect", UserLogin.Email)
		return nil, errors.Wrap(errors.New("Password isn't like password from database"), "Password incorrect")
	}

	var tokenString string
	tokenString, error := jwt.GenToken(types.UserFieldInToken{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email}, s.conf.Jwt.Duration)

	if error != nil {
		s.logger.Errorf("Can not gen token", error)
		return nil, errors.Wrap(error, "Can't gen token")
	}
	s.logger.Infof("Login completed ", user.Email)
	return &types.UserResponseSignUp{
		Name:  user.Name,
		Email: user.Email,
		Token: tokenString}, nil
}

// Get basic info for a user
func (s *Service) FindUserById(ctx context.Context, id string) (*types.UserResGetInfo, error) {
	var user *types.UserResGetInfo

	//check id correct
	if !bson.IsObjectIdHex(id) {
		s.logger.Errorf("Id user incorrect,it isn't ObjectIdHex")
		return nil, errors.New("Id incorrect to find the given user from database, it isn't ObjectIdHex")
	}

	user, err := s.repo.FindByID(ctx, id)

	if err != nil {
		s.logger.Errorf("Not found id user", err)
		return nil, errors.Wrap(err, "Failed to find id user from database")
	}
	s.logger.Infof("Find id completed ", id)
	return user, nil
}

// Post update info for a user
func (s *Service) UpdateUserByID(ctx context.Context, user types.User) error {

	err := s.repo.UpdateUserByID(ctx, user)

	if err != nil {
		s.logger.Errorf("failed when update user by id", err)
		return err
	}
	s.logger.Infof("updated user is completed ")
	return err
}

// Get list users by page
func (s *Service) GetListUsers(ctx context.Context, page, size string) (*types.GetListUsersResponse, error) {

	var pagingNSorting types.PagingNSorting

	if err := pagingNSorting.Init(page, size); err != nil {
		s.logger.Errorf("Failed url parameters when get list users", err)
		return nil, errors.Wrap(err, "Failed url parameters when get list users")
	}

	var listUsersResponse types.GetListUsersResponse

	numberUsers, err := s.repo.CountUser(ctx)
	if err != nil {
		s.logger.Errorf("Failed when get number users", err)
		return nil, errors.Wrap(err, "Failed when get number users")
	}

	listUsersResponse.CurrentPage = pagingNSorting.Page
	listUsersResponse.MaxItemsPerPage = pagingNSorting.Size
	listUsersResponse.TotalItems = numberUsers
	listUsersResponse.TotalPages = int(numberUsers / pagingNSorting.Size)
	// ex: total: 5, size: 2 => 3 page
	if numberUsers%pagingNSorting.Size != 0 {
		listUsersResponse.TotalPages += 1
	}

	if pagingNSorting.Size > numberUsers {
		listUsersResponse.MaxItemsPerPage = numberUsers
	}

	listUsers, err := s.repo.GetListUsers(ctx, pagingNSorting)

	if err != nil {
		s.logger.Errorf("Failed when get list users by page", err)
		return nil, errors.Wrap(err, "Failed when get list users by page")
	}

	listUsersResponse.ListUsers = append(listUsersResponse.ListUsers, listUsers...)
	s.logger.Infof("get list users by page is completed, page: ", pagingNSorting)

	return &listUsersResponse, nil
}

// get list user liked
func (s *Service) listLiked(ctx context.Context, userID string) ([]types.UserResGetInfo, error) {
	list, err := s.repo.GetListlikedInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.convertPointerArrayToArray(list), err
}

// get list user matched
func (s *Service) listMatched(ctx context.Context, userID string) ([]types.UserResGetInfo, error) {
	list, err := s.repo.GetListMatchedInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.convertPointerArrayToArray(list), err
}

// get list matched or liked
func (s *Service) GetMatchedUsersByID(ctx context.Context, idUser, matchedParameter string) ([]types.UserResGetInfo, error) {

	matched, err := strconv.ParseBool(matchedParameter)

	if err != nil {
		s.logger.Errorf("Failed url parameters when get list matched or like", err)
		return nil, errors.Wrap(err, "Failed url parameters when get list users")
	}

	if !bson.IsObjectIdHex(idUser) {
		s.logger.Errorf("Id user incorrect,it isn't ObjectIdHex")
		return nil, errors.New("Id user incorrect to find list liked from database, it isn't ObjectIdHex")
	}

	if matched {
		list, err := s.listMatched(ctx, idUser)
		s.logger.Infof("Get list matched completed", idUser)
		return list, err
	}

	list, err := s.listLiked(ctx, idUser)
	s.logger.Infof("Get list liked completed", idUser)
	return list, err
}

// convert []*types.UserResGetInfo to []types.UserResGetInfo - if empty return []
func (s *Service) convertPointerArrayToArray(list []*types.UserResGetInfo) []types.UserResGetInfo {

	listUsers := []types.UserResGetInfo{}
	for _, user := range list {
		listUsers = append(listUsers, *user)
	}
	return listUsers
}

// helps Enable/Disable account
func (s *Service) DisableUserByID(ctx context.Context, idUser, disableParameter string) error {

	disable, err := strconv.ParseBool(disableParameter)

	if err != nil {
		s.logger.Errorf("Failed url parameters when get list matched or like", err)
		return errors.Wrap(err, "Failed url parameters when get list users")
	}

	if !bson.IsObjectIdHex(idUser) {
		s.logger.Errorf("Id user incorrect,it isn't ObjectIdHex")
		return errors.New("Id user incorrect to find list liked from database, it isn't ObjectIdHex")
	}

	if disable {
		if err := s.repo.DisableUserByID(ctx, idUser, disable); err != nil {
			s.logger.Errorf("Disable failed", err)
		}
		s.logger.Infof("Disable completed", idUser)

		return err
	}

	if err = s.repo.DisableUserByID(ctx, idUser, disable); err != nil {
		s.logger.Errorf("Enable failed", err)
	}
	s.logger.Infof("Enable completed", idUser)
	return err
}
