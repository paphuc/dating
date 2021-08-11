package userservices

import (
	"context"
	"testing"
	"time"

	"dating/internal/app/api/types"
	"dating/internal/app/config"
	"dating/internal/pkg/glog"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepositoryMock struct {
	mock.Mock
}

func (mock *RepositoryMock) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.User), args.Error(1)
}

func (mock *RepositoryMock) CountUser(ctx context.Context, ps types.PagingNSorting) (int64, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(int64), args.Error(1)
}

func (mock *RepositoryMock) FindByID(ctx context.Context, id string) (*types.UserResGetInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.UserResGetInfo), args.Error(1)
}

func (mock *RepositoryMock) Insert(ctx context.Context, User types.User) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *RepositoryMock) UpdateUserByID(ctx context.Context, User types.User) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *RepositoryMock) GetListUsers(ctx context.Context, ps types.PagingNSorting) ([]*types.UserResGetInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*types.UserResGetInfo), args.Error(1)
}

func (mock *RepositoryMock) GetListlikedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*types.UserResGetInfo), args.Error(1)
}

func (mock *RepositoryMock) GetListMatchedInfo(ctx context.Context, idUser string) ([]*types.UserResGetInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*types.UserResGetInfo), args.Error(1)
}

func (mock *RepositoryMock) DisableUserByID(ctx context.Context, idUser string, disable bool) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *RepositoryMock) GetListUsersAvailable(ctx context.Context, ignoreIds []primitive.ObjectID, ps types.PagingNSorting) ([]*types.UserResGetInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*types.UserResGetInfo), args.Error(1)
}

func (mock *RepositoryMock) IgnoreIdUsers(ctx context.Context, id string) ([]primitive.ObjectID, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]primitive.ObjectID), args.Error(1)
}

func (mock *RepositoryMock) CountUserUsersAvailable(ctx context.Context, ignoreIds []primitive.ObjectID, ps types.PagingNSorting) (int64, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(int64), args.Error(1)
}

func userMock() (*types.User, error) {
	id, errID := primitive.ObjectIDFromHex("60e3f9a7e1ab4c3dfc8fe4c1")
	if errID != nil {
		return nil, errID
	}

	timeStr := "2021-07-11T17:00:00.000Z"
	timeDate, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}
	timeCreateStr := "2021-07-06T06:35:19.660Z"
	timeCreate, err := time.Parse(time.RFC3339, timeCreateStr)
	if err != nil {
		return nil, err
	}
	timeUpdateStr := "2021-08-02T05:01:56.284Z"
	timeUpdate, err := time.Parse(time.RFC3339, timeUpdateStr)
	if err != nil {
		return nil, err
	}

	user := &types.User{
		ID:           id,
		Name:         "Nhân Đinh Đẹp Trai",
		Email:        "tphuc@gmail.com",
		Birthday:     timeDate,
		Relationship: "FA mai mai luôn",
		LookingFor:   "Top",
		Password:     "$2a$10$rrXG4Uas1jltU.gPp.v2WexzAnwn5ZfOsl4rC2USOtJT5u5YnBnbW",
		Media: []string{
			"https://res.cloudinary.com/dng1knkia/image/upload/v1627879820/107801703_283005589704336_1892695093188363494_n_hdbpzi.jpg",
			"https://s1.img.yan.vn/YanNews/2167221/201603/20160310-124800-1_600x450.jpg",
			"https://i.pinimg.com/736x/ae/9d/b1/ae9db172c2735de851c45b82d7a988f4.jpg",
			"https://s1.img.yan.vn/YanNews/2167221/201603/20160310-124800-1_600x450.jpg",
			"https://i.pinimg.com/originals/57/64/32/576432ab92e270631eaab49f5a78f355.png",
			"https://s1.img.yan.vn/YanNews/2167221/201603/20160310-124800-1_600x450.jpg",
			"http://res.cloudinary.com/dng1knkia/image/upload/v1627467790/omc4580dkpgpenhbxex9.jpg",
			"http://res.cloudinary.com/dng1knkia/image/upload/v1627543866/paxxxj8eegek9sud6n4i.jpg",
			"http://res.cloudinary.com/dng1knkia/image/upload/v1627728351/ttmsrjr3i6nhh5asuylz.png",
		},
		Gender:  "Male",
		Sex:     "Male",
		Country: "Dong Thap",
		Hobby: []string{
			"Work", "Learn", "Kiss", "Sleepingggggggggggggg", "Pate heo",
		},
		Disable:  false,
		About:    "Động vật có vú, nhỏ nhắn và chuyên ăn thịt, sống chung với loài người, được   nuôi để săn vật gây hại hoặc làm thú nuôi cùng với chó nhà. ",
		CreateAt: timeCreate,
		UpdateAt: timeUpdate,
	}
	return user, nil
}

func userInfoMock() (*types.UserResGetInfo, error) {
	id, errID := primitive.ObjectIDFromHex("60e3f9a7e1ab4c3dfc8fe4c1")
	if errID != nil {
		return nil, errID
	}

	timeStr := "2021-07-11T17:00:00.000Z"
	timeDate, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}
	timeCreateStr := "2021-07-06T06:35:19.660Z"
	timeCreate, err := time.Parse(time.RFC3339, timeCreateStr)
	if err != nil {
		return nil, err
	}
	timeUpdateStr := "2021-08-02T05:01:56.284Z"
	timeUpdate, err := time.Parse(time.RFC3339, timeUpdateStr)
	if err != nil {
		return nil, err
	}

	user := &types.UserResGetInfo{
		ID:           id,
		Name:         "Nhân Đinh Đẹp Trai",
		Email:        "tphuc@gmail.com",
		Birthday:     timeDate,
		Relationship: "FA mai mai luôn",
		LookingFor:   "Top",
		Media: []string{
			"https://res.cloudinary.com/dng1knkia/image/upload/v1627879820/107801703_283005589704336_1892695093188363494_n_hdbpzi.jpg",
			"https://s1.img.yan.vn/YanNews/2167221/201603/20160310-124800-1_600x450.jpg",
			"https://i.pinimg.com/736x/ae/9d/b1/ae9db172c2735de851c45b82d7a988f4.jpg",
			"https://s1.img.yan.vn/YanNews/2167221/201603/20160310-124800-1_600x450.jpg",
			"https://i.pinimg.com/originals/57/64/32/576432ab92e270631eaab49f5a78f355.png",
			"https://s1.img.yan.vn/YanNews/2167221/201603/20160310-124800-1_600x450.jpg",
			"http://res.cloudinary.com/dng1knkia/image/upload/v1627467790/omc4580dkpgpenhbxex9.jpg",
			"http://res.cloudinary.com/dng1knkia/image/upload/v1627543866/paxxxj8eegek9sud6n4i.jpg",
			"http://res.cloudinary.com/dng1knkia/image/upload/v1627728351/ttmsrjr3i6nhh5asuylz.png",
		},
		Gender:  "Male",
		Sex:     "Male",
		Country: "Dong Thap",
		Hobby: []string{
			"Work", "Learn", "Kiss", "Sleepingggggggggggggg", "Pate heo",
		},
		Disable:  false,
		About:    "Động vật có vú, nhỏ nhắn và chuyên ăn thịt, sống chung với loài người, được   nuôi để săn vật gây hại hoặc làm thú nuôi cùng với chó nhà. ",
		CreateAt: timeCreate,
		UpdateAt: timeUpdate,
	}
	return user, nil
}

func TestLogin(t *testing.T) {

	user, err := userMock()
	if err != nil {
		assert.Error(t, err)
		return
	}

	mockRepo := new(RepositoryMock)
	mockRepo.On("FindByEmail").Return(user, nil).Once()
	mockRepo.On("FindByEmail").Return(nil, errors.New("Not Find user by email"))

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)
	result, err := testService.Login(context.Background(), types.UserLogin{
		Email:    "tphuc@gmail.com",
		Password: "123456",
	})
	if err != nil {
		assert.Error(t, err)
		return
	}
	assert.Equal(t, "Nhân Đinh Đẹp Trai", result.Name)
	result, err = testService.Login(context.Background(), types.UserLogin{
		Email:    "tphuc@gmail.com",
		Password: "123456",
	})
	if err != nil {
		assert.Error(t, err)
		return
	}
	mockRepo.AssertExpectations(t)
	assert.Equal(t, nil, result)
}

func TestSignUp(t *testing.T) {
	mockRepo := new(RepositoryMock)

	user, err := userMock()
	if err != nil {
		assert.Error(t, err)
		return
	}

	mockRepo.On("FindByEmail").Return(nil, errors.New("Email email exits")).Once()
	mockRepo.On("Insert").Return(nil)

	mockRepo.On("FindByEmail").Return(user, nil)

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	result, err := testService.SignUp(context.Background(), types.UserSignUp{
		Email:    "tphuc@gmail.com",
		Password: "123456",
		Name:     "Nhân Đinh Đẹp Trai",
	})
	if err != nil {
		assert.Error(t, err)
		return
	}
	result2, err := testService.SignUp(context.Background(), types.UserSignUp{
		Email:    "tphuc@gmail.com",
		Password: "123456",
		Name:     "Nhân Đinh Đẹp Trai",
	})
	if err != nil {
		assert.Error(t, err)
		return
	}
	mockRepo.AssertExpectations(t)
	assert.Equal(t, "Nhân Đinh Đẹp Trai", result.Name)
	assert.Equal(t, "tphuc@gmail.com", result.Email)
	assert.Equal(t, nil, result2)
}

func TestFindUserById(t *testing.T) {

	userInfo, err := userInfoMock()
	if err != nil {
		assert.Error(t, err)
		return
	}

	mockRepo := new(RepositoryMock)

	mockRepo.On("FindByID").Return(userInfo, nil).Once()
	mockRepo.On("FindByID").Return(nil, errors.New("Not found id"))

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	result, err := testService.FindUserById(context.Background(), "tphuc@gmail.com")
	if err != nil {
		assert.Error(t, err)
		return
	}
	result2, err := testService.FindUserById(context.Background(), "tphuc@gmail.com")
	if err != nil {
		assert.Error(t, err)
		return
	}
	mockRepo.AssertExpectations(t)
	assert.Equal(t, "Nhân Đinh Đẹp Trai", result.Name)
	assert.Equal(t, "tphuc@gmail.com", result.Email)
	assert.Equal(t, nil, result2)

}

func TestUpdateUserByID(t *testing.T) {

	user, err := userMock()
	if err != nil {
		assert.Error(t, err)
		return
	}

	mockRepo := new(RepositoryMock)

	mockRepo.On("UpdateUserByID").Return(nil).Once()
	mockRepo.On("UpdateUserByID").Return(errors.New("Failed when update user"))
	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	error := testService.UpdateUserByID(context.Background(), *user)
	error2 := testService.UpdateUserByID(context.Background(), *user)

	mockRepo.AssertExpectations(t)
	assert.Equal(t, nil, error)
	assert.NotEqual(t, nil, error2)
}
func TestGetListUsers(t *testing.T) {
	userInfo, _ := userInfoMock()

	mockRepo := new(RepositoryMock)

	mockRepo.On("CountUser").Return(int64(8), nil)
	mockRepo.On("GetListUsers").Return([]*types.UserResGetInfo{
		userInfo, userInfo,
		userInfo, userInfo,
	}, nil).Once()
	mockRepo.On("GetListUsers").Return(nil, errors.New("Failed when get list"))
	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	result, err := testService.GetListUsers(context.Background(), "1", "4", "18", "22", "Male")
	if err != nil {
		assert.Error(t, err)
		return
	}
	_, err = testService.GetListUsers(context.Background(), "1", "4", "18", "22", "Male")
	if err != nil {
		assert.Error(t, err)
		return
	}
	mockRepo.AssertExpectations(t)
	assert.Equal(t, 8, result.TotalItems)
	assert.Equal(t, 4, len(result.Content))
	assert.Equal(t, []string{"Male"}, result.Filter.Gender)
}

func TestGetMatchedUsersByID(t *testing.T) {

	userInfo, _ := userInfoMock()

	mockRepo := new(RepositoryMock)

	mockRepo.On("GetListlikedInfo").Return([]*types.UserResGetInfo{
		userInfo, userInfo,
		userInfo, userInfo,
	}, nil)
	mockRepo.On("GetListMatchedInfo").Return([]*types.UserResGetInfo{
		userInfo, userInfo,
		userInfo, userInfo, userInfo,
	}, nil)

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	result1, err := testService.GetMatchedUsersByID(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", "true")
	if err != nil {
		assert.Error(t, err)
		return
	}
	result2, err := testService.GetMatchedUsersByID(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", "false")
	if err != nil {
		assert.Error(t, err)
		return
	}
	_, err = testService.GetMatchedUsersByID(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", "falses")
	if err != nil {
		assert.Error(t, err)
		return
	}
	mockRepo.AssertExpectations(t)
	assert.Equal(t, 5, len(result1.Content))
	assert.Equal(t, 4, len(result2.Content))
}

func TestDisableUserByID(t *testing.T) {
	mockRepo := new(RepositoryMock)

	mockRepo.On("DisableUserByID").Return(nil).Twice()
	mockRepo.On("DisableUserByID").Return(errors.New("DisableUserByID Failed"))

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)
	errDisabled := testService.DisableUserByID(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", true)
	errEnabled := testService.DisableUserByID(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", false)
	errEnabled2 := testService.DisableUserByID(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", false)

	mockRepo.AssertExpectations(t)
	assert.Equal(t, nil, errDisabled)
	assert.Equal(t, nil, errEnabled)
	assert.Error(t, errEnabled2)
}
func TestGetListUsersAvailable(t *testing.T) {

	id1, err := primitive.ObjectIDFromHex("60e3b5d2e1ab4c388ce2d04a")
	if err != nil {
		assert.Error(t, err)
		return
	}
	id2, err := primitive.ObjectIDFromHex("60e3b5dbe1ab4c388ce2d04c")
	if err != nil {
		assert.Error(t, err)
		return
	}
	id3, err := primitive.ObjectIDFromHex("60e3f9a7e1ab4c3dfc8fe4c1")
	if err != nil {
		assert.Error(t, err)
		return
	}

	userInfo, _ := userInfoMock()

	mockRepo := new(RepositoryMock)

	mockRepo.On("IgnoreIdUsers").Return([]primitive.ObjectID{id1, id2, id3}, nil)
	mockRepo.On("CountUserUsersAvailable").Return(int64(2), nil)
	mockRepo.On("GetListUsersAvailable").Return([]*types.UserResGetInfo{
		userInfo, userInfo,
	}, nil)

	testService := NewService(
		&config.Configs{},
		&config.ErrorMessage{},
		mockRepo,
		glog.New(),
	)

	result, err := testService.GetListUsersAvailable(context.Background(), "60e3f9a7e1ab4c3dfc8fe4c1", "1", "4", "18", "21", "")
	if err != nil {
		assert.Error(t, err)
		return
	}

	mockRepo.AssertExpectations(t)
	assert.Equal(t, 2, result.TotalItems)
	assert.Equal(t, 2, len(result.Content))
	assert.Equal(t, []string{
		"Male",
		"Female",
		"Both"}, result.Filter.Gender)
}
