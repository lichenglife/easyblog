package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// 引入 mock 包
	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserBiz 用于mock userv1.UserBiz ,实现userv1.UserBiz 接口
type MockUserBiz struct {
	mock.Mock
}

// Create 创建用户
func (m *MockUserBiz) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserInfo, error) {
	args := m.Called(ctx, req)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*model.UserInfo), nil
}

// GetByID 根据 ID 获取用户
func (m *MockUserBiz) GetUserByID(ctx context.Context, userID string) (*model.UserInfo, error) {
	args := m.Called(ctx, userID)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*model.UserInfo), nil
}

// GetByUsername 根据用户名获取用户
func (m *MockUserBiz) GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error) {
	args := m.Called(ctx, username)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*model.UserInfo), nil
}

// Update 更新用户
func (m *MockUserBiz) UpdateUser(ctx context.Context, user *model.UpdateUser) error {
	args := m.Called(ctx, user)
	if args.Get(0) != nil {
		return nil
	}
	return args.Error(0)
}

// Delete 删除用户
func (m *MockUserBiz) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// List 获取用户列表
func (m *MockUserBiz) ListUsers(ctx context.Context, page, pageSize int) (*model.ListUserResponse, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*model.ListUserResponse), nil
}

// UserLogin 用户登录
func (m *MockUserBiz) UserLogin(ctx context.Context, user model.UserLoginRequest) (*model.UserInfo, error) {
	args := m.Called(ctx, user)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*model.UserInfo), nil
}

// ChangePassword 更新用户密码
func (m *MockUserBiz) ChangePassword(ctx context.Context, userID string, user model.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, user)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// MockBiz 用于mock userv1.UserBiz 接口
// 包含userv1.UserBiz 接口的所有方法
//
//	userHandlerTest 构建userHandler 测试结构体
type usertHandlerTest struct {
	// mock biz
	MockUserBiz *MockUserBiz
	mockBiz     *MockBiz
	handler     UserHandler
	// mock router
	router *gin.Engine
}
type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func setUpUserHandlerTest() *usertHandlerTest {
	// mock biz
	mockUserBiz := new(MockUserBiz)
	mockBiz := &MockBiz{UserBiz: mockUserBiz}
	// mock router
	router := gin.New() // mock handler
	handler := NewUserHandler(nil, mockBiz)
	return &usertHandlerTest{
		MockUserBiz: mockUserBiz,
		mockBiz:     mockBiz,
		handler:     handler,
		router:      router,
	}
}

// CreteUser 创建用户
func Test_CreateUser_SUCEESS(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.POST("/user", userHandlerTest.handler.CreateUser)

	// 3、准备测试数据
	req := &model.CreateUserRequest{
		Username: "zhangsan",
		Password: "Aa123456",
		Nickname: "zhangsang",
		Phone:    "15856695636",
		Email:    "123456@qq.com",
	}
	bodybytes, _ := json.Marshal(req)
	// 4、构造http请求和响应记录器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("POST", "/posts", strings.NewReader(string(bodybytes)))
	reqHttp.Header.Set("Content-Type", "application/json")

	// 5、设置gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = reqHttp
	c.Set("userID", "1")

	// 6、mock userbiz层返回
	userHandlerTest.MockUserBiz.On("CreateUser", mock.Anything, mock.Anything).Return(&model.UserInfo{Username: "zhangsan", Nickname: "zhangsang",
		Phone: "15856695636",
		Email: "123456@qq.com",
	}, nil).Once()

	// 7、执行handler
	userHandlerTest.handler.CreateUser(c)

	// 8、设置断言
	assert.Equal(t, http.StatusOK, w.Code)
	userHandlerTest.MockUserBiz.AssertExpectations(t)
}

func Test_CreateUser_Failed(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.POST("/user", userHandlerTest.handler.CreateUser)

	// 3、准备测试数据
	req := &model.CreateUserRequest{
		Username: "zhangsan",
		Password: "Aa123456",
		Nickname: "zhangsang",
		Phone:    "15856695636",
		Email:    "123456@qq.com",
	}
	bodybytes, _ := json.Marshal(req)
	// 4、构造http请求和响应记录器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("POST", "/posts", strings.NewReader(string(bodybytes)))
	reqHttp.Header.Set("Content-Type", "application/json")

	// 5、设置gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = reqHttp
	c.Set("userID", "1")

	// 6、mock userbiz层返回
	userHandlerTest.MockUserBiz.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errno.ErrUserAlreadyExist).Once()

	// 7、执行handler
	userHandlerTest.handler.CreateUser(c)

	// 8、设置断言
	assert.Equal(t, http.StatusConflict, w.Code)
	userHandlerTest.MockUserBiz.AssertExpectations(t)
}

func Test_CreateUser_ParamsISInValid(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.POST("/user", userHandlerTest.handler.CreateUser)

	// 3、准备测试数据
	req := &model.CreateUserRequest{
		Username: "zhangsan",
	}
	bodybytes, _ := json.Marshal(req)
	// 4、构造http请求和响应记录器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("POST", "/posts", strings.NewReader(string(bodybytes)))
	reqHttp.Header.Set("Content-Type", "application/json")

	// 5、设置gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = reqHttp
	c.Set("userID", "1")
	// 7、执行handler
	userHandlerTest.handler.CreateUser(c)

	// 8、设置断言
	assert.Equal(t, http.StatusOK, w.Code)
	userHandlerTest.MockUserBiz.AssertExpectations(t)
}

// ChangePassword 修改密码
// ResetPassword 重置密码
func Test_ResetPassword(t *testing.T) {

}

// UserInfo 获取用户信息
func Test_GetUserByID_SUCEESS(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.GET("/user/:userID", userHandlerTest.handler.GetUserByID)

	// 3、准备测试数据

	//4、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("GET", "/user/userID-1", nil)
	reqHttp.Header.Set("Content-Type", "application/json")

	// 5、设置gin context
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "userID-1")
	// 6、mockUserbiz
	userHandlerTest.MockUserBiz.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.UserInfo{Username: "zhangsan", Nickname: "zhangsang",

		Phone:  "15856695636",
		Email:  "123456@qq.com",
		UserID: "userID-1",
	}, nil).Once()

	// 7、执行测试handler
	userHandlerTest.handler.GetUserByID(c)

	// 8、设置断言
	assert.Equal(t, http.StatusOK, w.Code)
	userHandlerTest.MockUserBiz.AssertExpectations(t)
}

func Test_GetUserByID_Failed(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.GET("/user/:userID", userHandlerTest.handler.GetUserByID)

	// 3、准备测试数据

	//4、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("GET", "/user/userID-1", nil)
	reqHttp.Header.Set("Content-Type", "application/json")

	// 5、设置gin context
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "userID-1")
	// 6、mockUserbiz
	userHandlerTest.MockUserBiz.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, errno.ErrNotFound).Once()

	// 7、执行测试handler
	userHandlerTest.handler.GetUserByID(c)

	// 8、设置断言
	assert.Equal(t, http.StatusNotFound, w.Code)
	userHandlerTest.MockUserBiz.AssertExpectations(t)
}

// UserLogin 用户登录
func Test_UserLogin(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.POST("/user/login", userHandlerTest.handler.UserLogin)

	// 3、准备测试数据
	req := &model.UserLoginRequest{
		Username: "zhangsan",
		Password: "Aa123456",
	}
	bodybytes, _ := json.Marshal(req)
	//4、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("POST", "/user/login", strings.NewReader(string(bodybytes)))
	reqHttp.Header.Set("Content-Type", "application/json")
	// 5、设置gin context
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = reqHttp
	// 6、mockUserbiz
	userHandlerTest.MockUserBiz.On("UserLogin", mock.Anything, mock.Anything).Return(&model.UserInfo{Username: "zhangsan", Nickname: "zhangsang",
		Phone:  "15856695636",
		Email:  "123456@qq.com",
		UserID: "userID-1",
	}, nil).Once()
	// 7、执行测试handler
	userHandlerTest.handler.UserLogin(ginCtx)
	// 8、设置断言
	assert.Equal(t, http.StatusOK, w.Code)
	// 9、解析响应体
	var resp CommonResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	// 将interface{} 序列化为 json字节
	dataBytes, _ := json.Marshal(resp.Data)

	var loginResp model.UserLoginResponse
	// 反序列化json字节为目标结构体
	err = json.Unmarshal(dataBytes, &loginResp)
	assert.NoError(t, err)
	// 10、断言
	assert.Equal(t, "zhangsan", loginResp.User.Username)
	assert.Equal(t, "zhangsang", loginResp.User.Nickname)
	assert.Equal(t, "15856695636", loginResp.User.Phone)
	assert.Equal(t, "123456@qq.com", loginResp.User.Email)
	assert.Equal(t, "userID-1", loginResp.User.UserID)

	userHandlerTest.MockUserBiz.AssertExpectations(t)
}

// UserLogout 用户登出
func Test_UserLogout(t *testing.T) {

}

// UserInfo 获取用户信息
func Test_UserInfo(t *testing.T) {

}

// ListUsers 获取用户列表
func Test_ListUsers(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.GET("/users", userHandlerTest.handler.ListUsers)
	// 3、准备测试数据
	userList := []model.UserInfo{
		{Username: "zhangsan",
			Nickname: "zhangsang",
			Phone:    "15856695636",
			Email:    "123456@qq.com",
			UserID:   "userID-1",
		},
	}
	listUserResp := &model.ListUserResponse{
		TotalCount: 2,
		User:       userList,
		HasMore:    false,
	}
	//4、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("GET", "/users", nil)
	reqHttp.Header.Set("Content-Type", "application/json")
	// 5、设置gin context
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = reqHttp
	// 6、mockUserbiz
	userHandlerTest.MockUserBiz.On("ListUsers", mock.Anything, mock.Anything, mock.Anything).Return(listUserResp, nil).Once()
	// 7、执行测试handler
	userHandlerTest.handler.ListUsers(ginCtx)
	// 8、设置断言
	assert.Equal(t, http.StatusOK, w.Code)
	// 9、解析响应体
	var resp CommonResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	dataBytes, _ := json.Marshal(resp.Data)
	var listResp model.ListUserResponse
	err = json.Unmarshal(dataBytes, &listResp)
	assert.NoError(t, err)
	// 10、断言
	assert.Equal(t, resp.Code, 0)
	assert.Equal(t, resp.Message, "OK")
	assert.Equal(t, 1, len(listResp.User))
}

func Test_ListUsers_Failed(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.GET("/users", userHandlerTest.handler.ListUsers)
	// 3、准备测试数据
	//4、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("GET", "/users", nil)
	reqHttp.Header.Set("Content-Type", "application/json")
	// 5、设置gin context
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = reqHttp
	// 6、mockUserbiz
	userHandlerTest.MockUserBiz.On("ListUsers", mock.Anything, mock.Anything, mock.Anything).Return(nil, errno.ErrDatabase).Once()
	// 7、执行测试handler
	userHandlerTest.handler.ListUsers(ginCtx)
	// 8、设置断言
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// 保证mock函数全部正确执行
	userHandlerTest.MockUserBiz.AssertExpectations(t)
	// 9、解析响应体
}

// GetUserByID 根据 ID 获取用户
func Test_GetUserByID(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.GET("/user/:userID", userHandlerTest.handler.GetUserByID)

	// 3、准备测试数据
	userID := "userID-1"
	//4、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("GET", "/user/"+userID, nil)
	reqHttp.Header.Set("Content-Type", "application/json")
	// 5、设置gin context
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = reqHttp
	// 6、mockUserbiz
	userHandlerTest.MockUserBiz.On("GetUserByID", mock.Anything, mock.Anything, mock.Anything).Return(&model.UserInfo{Username: "zhangsan", Nickname: "zhangsang",
		Phone:  "15856695636",
		Email:  "123456@qq.com",
		UserID: "userID-1",
	}, nil).Once()
	// 7、执行测试handler
	userHandlerTest.handler.GetUserByID(ginCtx)

	// 8、设置断言
	assert.Equal(t, http.StatusOK, w.Code)
	userHandlerTest.MockUserBiz.AssertExpectations(t)
	// 解析http response body
}

// UpdateUser 更新用户
func Test_UpdateUser(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.PUT("/user/:userID", userHandlerTest.handler.UpdateUser)

	// 3、准备测试数据
	userID := "userID-1"
	updateUserReq := &model.UpdateUser{
		Nickname: "zhangsang",
		Phone:    "15856695636",
		Email:    "123456@qq.com",
	}
	//4、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqData, _ := json.Marshal(updateUserReq)
	reqHttp, _ := http.NewRequest("PUT", "/user/"+userID, strings.NewReader(string(reqData)))
	reqHttp.Header.Set("ContentType", "application/json")

	// 5、设置gin context
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = reqHttp
	//ginCtx.Params = gin.Params{gin.Param{Key: "userID", Value: userID}}
	// 6、mockUserbiz
	userHandlerTest.MockUserBiz.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// 7、执行测试handler
	userHandlerTest.handler.UpdateUser(ginCtx)
	// 8、设置断言
	assert.Equal(t, http.StatusOK, w.Code)
}

// DeleteUser 删除用户  全流程测试
func Test_DeleteUser_Full_Process(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2. 注册中间件，设置 context 的值
	userHandlerTest.router.Use(func(c *gin.Context) {
		c.Set("username", "root")
		c.Set("userID", "root")
		c.Next()
	})
	// 3、注册路由
	userHandlerTest.router.DELETE("/users/:userID", userHandlerTest.handler.DeleteUser)
	// 4、mock
	userHandlerTest.MockUserBiz.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	userID := "userID-1"
	//5、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("DELETE", "/users/"+userID, nil)
	reqHttp.Header.Set("ContentType", "application/json")
	// 5、gin自动分发路由

	userHandlerTest.router.ServeHTTP(w, reqHttp)
	// 6、断言
	assert.Equal(t, http.StatusOK, w.Code)
	var resp CommonResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, resp.Code, 0)
	assert.Equal(t, resp.Message, "OK")
}

func Test_DeleteUser_Function(t *testing.T) {
	//  1、初始化 handler 及mock
	userHandlerTest := setUpUserHandlerTest()
	// 2、注册路由
	userHandlerTest.router.DELETE("/users/:userID", userHandlerTest.handler.DeleteUser)
	// 3、准备测试数据
	userID := "userID-1"
	//4、构建http请求以及响应器
	w := httptest.NewRecorder()
	reqHttp := httptest.NewRequest("DELETE", "/users/"+userID, nil)
	reqHttp.Header.Set("ContentType", "application/json")
	// 5、设置gin context
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Set("username", "root")
	ginCtx.Request = reqHttp
	// 6、mockUserbiz
	userHandlerTest.MockUserBiz.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	// 7、执行测试handler
	userHandlerTest.handler.DeleteUser(ginCtx)

	// 8、设置断言
	assert.Equal(t, http.StatusOK, w.Code)
}
