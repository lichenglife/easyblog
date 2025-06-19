package biz

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/utils/authn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// biz层单元测试

// biz 依赖于 store 层，所以需要对 store 层进行 mock
// 1. 直接使用 mock 库生成 mock 对象
// 2. 使用 gomock 库生成 mock 对象

// 1、定义UserMockStore 结构体，实现userStore全部接口
type MockUserStore struct {
	mock.Mock
}

// MockPostStore 模拟 PostStore 接口
type MockPostStore struct {
	mock.Mock
}

// MockPostStore 实现store.PostStore
func (m *MockPostStore) GetByUserID(ctx context.Context, userID string, page, pageSize int) (int64, []*model.Post, error) {
	args := m.Called(ctx, userID, page, pageSize)
	return args.Get(0).(int64), args.Get(1).([]*model.Post), args.Error(2)
}

// 其他方法返回默认值
func (m *MockPostStore) Create(ctx context.Context, post *model.Post) error { return nil }
func (m *MockPostStore) Delete(ctx context.Context, postID string) error    { return nil }
func (m *MockPostStore) GetByPostID(ctx context.Context, postID string) (*model.Post, error) {
	return nil, nil
}
func (m *MockPostStore) List(ctx context.Context, page, pageSize int) (int64, []*model.Post, error) {
	return 0, nil, nil
}
func (m *MockPostStore) Update(ctx context.Context, post *model.Post) error        { return nil }
func (m *MockPostStore) GetByID(ctx context.Context, id uint) (*model.Post, error) { return nil, nil }

// Create 创建用户
func (m *MockUserStore) Create(ctx context.Context, user *model.User) error {
	// 调用mock的方法
	args := m.Called(ctx, user)
	return args.Error(0)
}

// GetByID 根据 ID 获取用户
func (m *MockUserStore) GetByID(ctx context.Context, userID string) (*model.User, error) {
	// 调用mock的方法
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// GetByUsername 根据用户名获取用户
func (m *MockUserStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	// 调用mock的方法
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// Update 更新用户
func (m *MockUserStore) Update(ctx context.Context, user *model.User) error {
	// 调用mock的方法
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Delete 删除用户
func (m *MockUserStore) Delete(ctx context.Context, userID string) error {
	// 调用mock的方法
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// List 获取用户列表
func (m *MockUserStore) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	// 调用mock的方法
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

// MockStore 对IStore层进行mock, 实现IStore的接口
type MockStore struct {
	mock.Mock
	UserStore store.UserStore
	PostStore store.PostStore
}

// User 实现IStore的User方法
func (m *MockStore) User() store.UserStore {
	return m.UserStore
}

func (m *MockStore) Post() store.PostStore {
	return m.PostStore
}
func (m *MockStore) Close() error {
	return nil
}

// 执行单元测试

// Create 创建用户
func Test_CreateUser(t *testing.T) {
	// mock UserStore
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(mockStore)

	// mock request
	req := &model.CreateUserRequest{
		Username: "zjamhsa",
		Password: "zajmsah",
		Nickname: "zansag",
		Email:    "70985555@qq.com",
		Phone:    "4455222",
	}

	// mock userstore执行查询 GetByUsername
	mockUserStore.On("GetByUsername", mock.Anything, "zjamhsa").Return(nil, nil).Once()
	// mock  userstore被执行，返回nil
	mockUserStore.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()

	// 执行被测试函数
	user, err := biz.CreateUser(context.Background(), req)

	// 断言
	assert.NoError(t, err)
	assert.Equal(t, req.Username, user.Username)
	mockUserStore.AssertExpectations(t)
}

func Test_CreateUser_Failed(t *testing.T) {
	// mock UserStore
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(mockStore)

	// mock request
	req := &model.CreateUserRequest{
		Username: "zjamhsa",
		Password: "zajmsah",
		Nickname: "zansag",
		Email:    "70985555@qq.com",
		Phone:    "4455222",
	}
	existuser := &model.User{
		UserID:   "user-1",
		Username: "zjamhsa",
		Password: "zajmsah",
		NickName: "zansag",
		Email:    "70985555@qq.com",
		Phone:    "4455222",
	}
	// mock userstore执行查询 GetByUsername
	mockUserStore.On("GetByUsername", mock.Anything, "zjamhsa").Return(existuser, errors.New("user already exists")).Once()
	// mock  userstore被执行，返回nil
	//mockUserStore.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()

	// 执行被测试函数
	_, err := biz.CreateUser(context.Background(), req)
	// 断言
	assert.Error(t, err)
	assert.Error(t, errno.ErrUserAlreadyExist)
	mockUserStore.AssertExpectations(t)
}

// GetByID 根据 ID 获取用户
func Test_GetUserByID(t *testing.T) {
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(mockStore)
	existuser := &model.User{
		UserID:   "user-1",
		Username: "zjamhsa",
		Password: "zajmsah",
		NickName: "zansag",
		Email:    "70985555@qq.com",
		Phone:    "4455222",
	}
	// mock userstore执行查询 GetByID
	mockUserStore.On("GetByID", mock.Anything, "user-1").Return(existuser, nil).Once()

	user, err := biz.GetUserByID(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, existuser.Username, user.Username)
	mockUserStore.AssertExpectations(t)
}

func Test_GetUserByID_Failed(t *testing.T) {
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(mockStore)

	// mock userstore执行查询 GetByID
	mockUserStore.On("GetByID", mock.Anything, "user-1").Return(nil, errno.ErrUserNotFound).Once()

	_, err := biz.GetUserByID(context.Background(), "user-1")

	assert.Error(t, err)
	assert.Equal(t, err, errno.ErrUserNotFound)
	mockUserStore.AssertExpectations(t)
}

// GetByUsername 根据用户名获取用户
func Test_GetUserByUsername(t *testing.T) {
	mockUserStore := new(MockUserStore)
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
		PostStore: mockPostStore, // mockPostStore 必须实现 PostStore 接口的所有方法
	}

	biz := NewUserBiz(mockStore)
	existuser := &model.User{
		UserID:   "user-1",
		Username: "zjamhsa",
		Password: "zajmsah",
		NickName: "zansag",
		Email:    "70985555@qq.com",
		Phone:    "4455222",
	}

	// mock userstore执行查询 GetByID
	mockUserStore.On("GetByUsername", mock.Anything, "zjamhsa").Return(existuser, nil).Once()
	// mock poststore执行查询 GetByUserID
	var posts []*model.Post
	mockPostStore.On("GetByUserID", mock.Anything, "user-1", 1, 10).Return(int64(1), posts, nil).Once()

	user, err := biz.GetUserByUsername(context.Background(), "zjamhsa")

	assert.NoError(t, err)
	assert.Equal(t, existuser.Username, user.Username)
	mockUserStore.AssertExpectations(t)
}

// Update 更新用户
func Test_UpdateUser(t *testing.T) {
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(mockStore)

	req := &model.UpdateUser{
		Nickname: "zhangsan",
		UserID:   "user-1",
		Email:    "1522552@qq.com",
	}

	user := &model.User{
		NickName: "zhangsan",
		UserID:   "user-1",
		Email:    "1522552@qq.com",
	}
	mockUserStore.On("GetByID", mock.Anything, "user-1").Return(user, nil).Once()

	mockUserStore.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	err := biz.UpdateUser(context.Background(), req)
	assert.NoError(t, err)
	mockUserStore.AssertExpectations(t)

}

// Delete 删除用户
func Test_DeleteUser(t *testing.T) {
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(mockStore)

	mockUserStore.On("GetByID", mock.Anything, "user-1").Return(&model.User{
		UserID: "user-1",
	}, nil).Once()
	mockUserStore.On("Delete", mock.Anything, "user-1").Return(nil).Once()
	err := biz.DeleteUser(context.Background(), "user-1")
	// 断言
	assert.NoError(t, err)
	mockUserStore.AssertExpectations(t)
}

func Test_ListUsers_SUCESS(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
		PostStore: mockPostStore,
	}
	biz := NewUserBiz(mockStore)
	// 准备测试数据
	users := []*model.User{
		{
			UserID:    "user-1",
			Username:  "zhangsan",
			NickName:  "张三",
			Email:     "zhangsan@example.com",
			Phone:     "13800138000",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    "user-2",
			Username:  "lisi",
			NickName:  "李四",
			Email:     "lisi@example.com",
			Phone:     "13800138001",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 设置 mock 期望
	mockUserStore.On("List", mock.Anything, 1, 10).Return(users, int64(2), nil).Once()

	// 为每个用户 mock GetByUserID 调用
	for _, user := range users {
		mockPostStore.On("GetByUserID", mock.Anything, user.UserID, 1, 10).Return(int64(2), []*model.Post{}, nil).Once()
	}

	// 执行测试
	resp, err := biz.ListUsers(context.Background(), 1, 10)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.TotalCount)
	assert.Equal(t, 2, len(resp.User))

	// 由于并发处理，我们需要检查返回的用户列表是否包含所有用户，而不是检查顺序
	userMap := make(map[string]model.UserInfo)
	for _, user := range resp.User {
		userMap[user.UserID] = user
	}

	// 验证每个用户的信息
	assert.Contains(t, userMap, "user-1")
	assert.Contains(t, userMap, "user-2")
	assert.Equal(t, "zhangsan", userMap["user-1"].Username)
	assert.Equal(t, "lisi", userMap["user-2"].Username)
	assert.Equal(t, int64(2), userMap["user-1"].BlogTotal)
	assert.Equal(t, int64(2), userMap["user-2"].BlogTotal)

	// 验证所有 mock 调用
	mockUserStore.AssertExpectations(t)
	mockPostStore.AssertExpectations(t)
}

func Test_ListUsers_Failed(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
		PostStore: mockPostStore,
	}
	biz := NewUserBiz(mockStore)
	// 设置 mock 期望
	mockUserStore.On("List", mock.Anything, 1, 10).Return(nil, int64(2), errno.ErrDatabase).Once()
	// 执行测试
	_, err := biz.ListUsers(context.Background(), 1, 10)
	assert.Error(t, err)
	assert.EqualError(t, err, errno.ErrDatabase.Error())
	// 验证所有 mock 调用
	mockUserStore.AssertExpectations(t)
}

func Test_ListUsers_BlogCount_Failed(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockUserStore := new(MockUserStore)
	mockStore := &MockStore{
		UserStore: mockUserStore,
		PostStore: mockPostStore,
	}
	biz := NewUserBiz(mockStore)

	users := []*model.User{
		{
			UserID:    "user-1",
			Username:  "zhangsan",
			NickName:  "张三",
			Email:     "zhangsan@example.com",
			Phone:     "13800138000",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	// 设置 mock 期望
	mockUserStore.On("List", mock.Anything, 1, 10).Return(users, int64(1), nil).Once()
	for _, user := range users {
		mockPostStore.On("GetByUserID", mock.Anything, user.UserID, 1, 10).Return(int64(2), []*model.Post{}, errno.ErrDatabase).Once()
	}
	// 执行测试
	_, err := biz.ListUsers(context.Background(), 1, 10)
	assert.Error(t, err)
	assert.EqualError(t, err, errno.ErrDatabase.Error())
	// 验证所有 mock 调用
	mockUserStore.AssertExpectations(t)
}

// UserLogin 用户登录
func Test_UserLogin(t *testing.T) {

	mockUserStore := new(MockUserStore)

	store := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(store)
	// 准备数据
	req := &model.UserLoginRequest{
		Username: "zhangsan",
		Password: "1231",
	}
	user := &model.User{
		Username: "zhangsan",
		Password: "1231",
	}
	encryptPwd, _ := authn.Encrypt(req.Password)
	user.Password = encryptPwd
	// 设置mock
	mockUserStore.On("GetByUsername", mock.Anything, mock.Anything).Return(user, nil).Once()
	// 执行测试
	response, err := biz.UserLogin(context.Background(), *req)
	// 设置断言
	assert.NoError(t, err)
	assert.Equal(t, "zhangsan", response.Username)

	mockUserStore.AssertExpectations(t)
}

func Test_UserLogin_Failed(t *testing.T) {

	mockUserStore := new(MockUserStore)

	store := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(store)
	// 准备数据
	req := &model.UserLoginRequest{
		Username: "zhangsan",
		Password: "1231",
	}
	user := &model.User{
		Username: "zhangsan",
		Password: "1231",
	}
	// 设置mock
	mockUserStore.On("GetByUsername", mock.Anything, mock.Anything).Return(user, nil).Once()
	// 执行测试
	_, err := biz.UserLogin(context.Background(), *req)
	// 设置断言
	assert.Error(t, err)
	assert.EqualError(t, err, errno.ErrPasswordIncorrect.Error())

	// 设置mock
	mockUserStore.On("GetByUsername", mock.Anything, mock.Anything).Return(nil, errno.ErrNotFound).Once()
	// 执行测试
	_, err = biz.UserLogin(context.Background(), *req)
	// 设置断言
	assert.Error(t, err)
	mockUserStore.AssertExpectations(t)
}

// ChangePassword 更新用户密码
func Test_ChangePassword(t *testing.T) {
	//  创建数据库层mock实例

	mockUserStore := new(MockUserStore)
	store := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(store)
	// 准备数据
	req := &model.ChangePasswordRequest{
		OldPassword: "121",
		NewPassword: "123",
	}
	oldPwd, _ := authn.Encrypt(req.OldPassword)
	// 设置mock
	mockUserStore.On("GetByID", mock.Anything, "zhangsan").Return(&model.User{
		UserID:   "zhangsan",
		Username: "zhangsan",
		Password: oldPwd,
	}, nil).Once()
	mockUserStore.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

	// 执行测试
	err := biz.ChangePassword(context.Background(), "zhangsan", *req)
	// 设置断言
	assert.NoError(t, err)
	mockUserStore.AssertExpectations(t)

}

func Test_ChangePassword_Failed(t *testing.T) {
	//  创建数据库层mock实例

	mockUserStore := new(MockUserStore)
	store := &MockStore{
		UserStore: mockUserStore,
	}
	biz := NewUserBiz(store)
	// 准备数据
	req := &model.ChangePasswordRequest{
		OldPassword: "121",
		NewPassword: "123",
	}
	// 设置mock
	mockUserStore.On("GetByID", mock.Anything, "zhangsan").Return(nil, errno.ErrDatabase).Once()
	//mockUserStore.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
	// 执行测试
	err := biz.ChangePassword(context.Background(), "zhangsan", *req)
	// 设置断言
	assert.Error(t, err)
	assert.EqualError(t, err, errno.ErrDatabase.Error())
	mockUserStore.AssertExpectations(t)
}
