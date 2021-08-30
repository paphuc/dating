package userhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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

type ServiceMock struct {
	mock.Mock
}

func (mock *ServiceMock) SignUp(ctx context.Context, UserSignUp types.UserSignUp) (*types.UserResponseSignUp, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.UserResponseSignUp), args.Error(1)
}

func (mock *ServiceMock) Login(ctx context.Context, UserLogin types.UserLogin) (*types.UserResponseSignUp, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.UserResponseSignUp), args.Error(1)
}

func (mock *ServiceMock) FindUserById(ctx context.Context, id string) (*types.UserResGetInfo, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.UserResGetInfo), args.Error(1)
}

func (mock *ServiceMock) UpdateUserByID(ctx context.Context, User types.User) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *ServiceMock) GetListUsers(ctx context.Context, page, size, minAge, maxAge, gender string) (*types.GetListUsersResponse, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.GetListUsersResponse), args.Error(1)
}

func (mock *ServiceMock) GetMatchedUsersByID(ctx context.Context, idUser, matchedParameter string) (types.ListUsersResponse, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(types.ListUsersResponse), args.Error(1)
}

func (mock *ServiceMock) DisableUserByID(ctx context.Context, idUser string, disable bool) error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *ServiceMock) GetListUsersAvailable(ctx context.Context, id, page, size, minAge, maxAge, gender string) (*types.GetListUsersResponse, error) {

	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*types.GetListUsersResponse), args.Error(1)
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

func TestSignUp(t *testing.T) {

	mockService := new(ServiceMock)
	mockService.On("SignUp").Return(&types.UserResponseSignUp{
		Name:  "Nhân Đinh Đẹp Trai",
		Email: "tphuc@gmail.com",
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2MGUzZjlhN2UxYWI0YzNkZmM4ZmU0YzEiLCJlbWFpbCI6InRwaHVjQGdtYWlsLmNvbSIsIm5hbWUiOiJOaMOibiDEkGluaCDEkOG6uXAgVHJhaSIsImV4cCI6MTYyODY0ODY1N30.1L3hlj7F4W2dI1oJ44CWO-0o3IzYJvaWeG27I9AJt3Q",
	}, nil).Once()
	mockService.On("SignUp").Return(nil, errors.New("err"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	bodyStruct := types.UserSignUp{
		Name:     "Nhân Đinh Đẹp Trai",
		Email:    "tphuc@gmail.com",
		Password: "12345678",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.SignUp(w, r)
	}))
	defer func() { ts.Close() }()

	body, _ := json.Marshal(bodyStruct)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	body_res, _ := io.ReadAll(res.Body)

	var body_mock *types.UserResponseSignUp
	json.Unmarshal([]byte(body_res), &body_mock)

	bodyStruct2 := types.UserSignUp{
		Name:     "Nhân Đinh Đẹp Trai",
		Email:    "tphuc@gmail.com",
		Password: "123456",
	}
	body, _ = json.Marshal(bodyStruct2)
	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	body, _ = json.Marshal(`"email": "tphuc@gmail.com","a": "b"`)
	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	assert.NoError(t, err)
	assert.Equal(t, "Nhân Đinh Đẹp Trai", body_mock.Name)
	body3, _ := json.Marshal(bodyStruct)
	req3, err := http.NewRequest("POST", ts.URL, bytes.NewReader(body3))
	assert.NoError(t, err)
	res3, err := http.DefaultClient.Do(req3)

	assert.NoError(t, err)
	assert.Equal(t, 409, res3.StatusCode)

}

func TestLogin(t *testing.T) {
	// login completed
	mockService := new(ServiceMock)
	mockService.On("Login").Return(&types.UserResponseSignUp{
		Name:  "Nhân Đinh Đẹp Trai",
		Email: "tphuc@gmail.com",
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2MGUzZjlhN2UxYWI0YzNkZmM4ZmU0YzEiLCJlbWFpbCI6InRwaHVjQGdtYWlsLmNvbSIsIm5hbWUiOiJOaMOibiDEkGluaCDEkOG6uXAgVHJhaSIsImV4cCI6MTYyODY0ODY1N30.1L3hlj7F4W2dI1oJ44CWO-0o3IzYJvaWeG27I9AJt3Q",
	}, nil).Once()
	mockService.On("Login").Return(nil, errors.New("IncorrectPasswordEmail"))

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	bodyStruct := types.UserSignUp{
		Name:     "Nhân Đinh Đẹp Trai",
		Email:    "tphuc@gmail.com",
		Password: "12345678",
	}
	body, _ := json.Marshal(bodyStruct)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.Login(w, r)
	}))
	defer func() { ts.Close() }()

	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	var body_mock *types.UserResponseSignUp
	body_res, _ := io.ReadAll(res.Body)

	json.Unmarshal([]byte(body_res), &body_mock)
	assert.Equal(t, "Nhân Đinh Đẹp Trai", body_mock.Name)

	bodyStructNotEmail := types.UserSignUp{
		Name:     "Nhân Đinh Đẹp Trai",
		Email:    "tphucgmail.com",
		Password: "12345678",
	}
	body, _ = json.Marshal(bodyStructNotEmail)
	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	body, _ = json.Marshal(`"email": "tphuc@gmail.com","a": "b"`)
	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	body, _ = json.Marshal(bodyStruct)
	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

}
func TestUpdateUserByID(t *testing.T) {
	mockService := new(ServiceMock)

	mockService.On("UpdateUserByID").Return(nil).Once()
	mockService.On("UpdateUserByID").Return(errors.New("StatusInternalServerError"))
	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.UpdateUserByID(w, r)
	}))
	defer func() { ts.Close() }()

	user, err := userMock()
	if err != nil {
		assert.Error(t, err)
		return
	}
	user.Email = "userexample.com"

	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	body, _ = json.Marshal(`"email": "tphuc@gmail.com","a": "b"`)
	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	user.Email = "tphuc@gmail.com"
	body, _ = json.Marshal(user)
	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	user.Email = "tphuc@gmail.com"
	body, _ = json.Marshal(user)
	req, err = http.NewRequest("POST", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetListUsers(t *testing.T) {
	userInfo, _ := userInfoMock()
	mockService := new(ServiceMock)

	resultMock := &types.GetListUsersResponse{
		Pagination: types.Pagination{
			TotalItems:      10,
			TotalPages:      3,
			CurrentPage:     1,
			MaxItemsPerPage: 4,
			Filter: types.Filter{
				AgeRange: types.AgeRange{
					Gte: time.Now(),
					Lt:  time.Now(),
				},
				Gender: []string{
					"Male",
					"Female",
					"Other",
				},
			},
		},
		ListUsersResponse: types.ListUsersResponse{
			Content: []*types.UserResGetInfo{
				userInfo, userInfo,
				userInfo, userInfo,
			},
		},
	}
	mockService.On("GetListUsers").Return(resultMock, nil).Once()
	mockService.On("GetListUsers").Return(nil, errors.New("Can't get")).Once()

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.GetListUsers(w, r)
	}))
	defer func() { ts.Close() }()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body_res, _ := io.ReadAll(res.Body)
	var body_mock *types.GetListUsersResponse
	json.Unmarshal([]byte(body_res), &body_mock)
	assert.Equal(t, 10, body_mock.TotalItems)

	req, err = http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, res.StatusCode)

}

func TestGetMatchedUsersByID(t *testing.T) {
	userInfo, _ := userInfoMock()
	mockService := new(ServiceMock)

	resultMock := types.ListUsersResponse{
		Content: []*types.UserResGetInfo{
			userInfo, userInfo,
			userInfo, userInfo,
		}}
	mockService.On("GetMatchedUsersByID").Return(resultMock, nil).Once()
	mockService.On("GetMatchedUsersByID").Return(types.ListUsersResponse{Content: []*types.UserResGetInfo{}}, errors.New("Can't get")).Once()

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.GetMatchedUsersByID(w, r)
	}))
	defer func() { ts.Close() }()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body_res, _ := io.ReadAll(res.Body)
	var body_mock *types.ListUsersResponse
	json.Unmarshal([]byte(body_res), &body_mock)
	assert.Equal(t, 4, len(body_mock.Content))

	req, err = http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

}

func TestDisableUsersByID(t *testing.T) {
	mockService := new(ServiceMock)

	mockService.On("DisableUserByID").Return(nil).Once()
	mockService.On("DisableUserByID").Return(errors.New("Can't Disable")).Once()

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.DisableUsersByID(w, r)
	}))
	defer func() { ts.Close() }()

	a := true
	disable := types.DisableBody{Disable: &a}

	body, _ := json.Marshal(`"disable": "tphuc@gmail.com"`)
	req, err := http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	body, _ = json.Marshal(types.DisableBody{})
	req, err = http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

	body, _ = json.Marshal(disable)
	req, err = http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	body, _ = json.Marshal(disable)
	req, err = http.NewRequest("PATCH", ts.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)

}

func TestGetListUsersAvailable(t *testing.T) {
	userInfo, _ := userInfoMock()
	mockService := new(ServiceMock)

	resultMock := &types.GetListUsersResponse{
		Pagination: types.Pagination{
			TotalItems:      10,
			TotalPages:      3,
			CurrentPage:     1,
			MaxItemsPerPage: 4,
			Filter: types.Filter{
				AgeRange: types.AgeRange{
					Gte: time.Now(),
					Lt:  time.Now(),
				},
				Gender: []string{
					"Male",
					"Female",
					"Other",
				},
			},
		},
		ListUsersResponse: types.ListUsersResponse{
			Content: []*types.UserResGetInfo{
				userInfo, userInfo,
				userInfo, userInfo,
			},
		},
	}
	mockService.On("GetListUsersAvailable").Return(resultMock, nil).Once()
	mockService.On("GetListUsersAvailable").Return(nil, errors.New("Can't get")).Once()

	testHandler := New(
		&config.Configs{},
		&config.ErrorMessage{},
		mockService,
		glog.New(),
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testHandler.GetListUsersAvailable(w, r)
	}))
	defer func() { ts.Close() }()

	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	body_res, _ := io.ReadAll(res.Body)
	var body_mock *types.GetListUsersResponse
	json.Unmarshal([]byte(body_res), &body_mock)
	assert.Equal(t, 10, body_mock.TotalItems)

	req, err = http.NewRequest("GET", ts.URL, nil)
	assert.NoError(t, err)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, res.StatusCode)
}
