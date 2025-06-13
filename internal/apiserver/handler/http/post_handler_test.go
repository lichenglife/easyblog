package handler

import (
	"context"
	"encoding/json"
	"errors"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bytes"

	"github.com/gin-gonic/gin"
	postv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/post"
	userv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/user"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostBiz 用于mock postv1.PostBiz
// 只需实现用到的方法
type MockPostBiz struct {
	mock.Mock
}

func (m *MockPostBiz) CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}
func (m *MockPostBiz) DeletePostByPostID(ctx context.Context, postID string) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}
func (m *MockPostBiz) GetPostByPostID(ctx context.Context, postID string) (*model.Post, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}
func (m *MockPostBiz) GetPostsByUserID(ctx context.Context, userID string, page, pageSize int) (*model.ListPostResponse, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ListPostResponse), args.Error(1)
}
func (m *MockPostBiz) ListPosts(ctx context.Context, page, pageSize int) (*model.ListPostResponse, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ListPostResponse), args.Error(1)
}
func (m *MockPostBiz) UpdatePost(ctx context.Context, req *model.UpdatePostRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// MockBiz 用于mock IBiz
// 只需实现 PostV1()
type MockBiz struct {
	mock.Mock
	PostBiz postv1.PostBiz
}

func (m *MockBiz) PostV1() postv1.PostBiz {
	return m.PostBiz
}

func (m *MockBiz) UserV1() userv1.UserBiz { return nil }

// postHandlerTest构建测试结构体
type postHandlerTest struct {
	// mockPostBiz
	mockPostBiz *MockPostBiz
	mockBiz     *MockBiz
	handler     PostHandler
	router      *gin.Engine
}

// setupPostHandlerTest
func setupPostHandlerTest() *postHandlerTest {
	mockPostBiz := new(MockPostBiz)
	mockBiz := &MockBiz{PostBiz: mockPostBiz}
	handler := NewPostHandler(nil, mockBiz)
	router := gin.New()
	return &postHandlerTest{
		mockPostBiz: mockPostBiz,
		mockBiz:     mockBiz,
		handler:     handler,
		router:      router,
	}
}

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TestPostHandler_CreatePost 测试创建帖子接口
func TestPostHandler_CreatePost(t *testing.T) {

	postHandlerTest := setupPostHandlerTest()
	postHandlerTest.router.POST("/posts", postHandlerTest.handler.CreatePost)
	// mock 请求数据
	reqBody := &model.CreatePostRequest{
		Content: "测试",
		Title:   "测试",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// 模拟响应写入
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// 5、mock  http 请求
	c.Request = httptest.NewRequest("POST", "/posts", strings.NewReader(string(bodyBytes)))
	c.Set("userID", "u1")
	// 6、mock postBiz业务层返回值
	postHandlerTest.mockPostBiz.On("CreatePost", mock.Anything, mock.AnythingOfType("*model.CreatePostRequest")).Return(&model.Post{PostID: "p1", UserID: "u1"}, nil).Once()
	// 7、执行测试函数
	postHandlerTest.handler.CreatePost(c)
	// 8、断言
	assert.Equal(t, http.StatusOK, w.Code)
	postHandlerTest.mockPostBiz.AssertExpectations(t)
}

// TestPostHandler_CreatePost_Fail 测试创建帖子失败
func TestPostHandler_CreatePost_Fail(t *testing.T) {
	mockPostBiz := new(MockPostBiz)
	mockBiz := &MockBiz{PostBiz: mockPostBiz}
	h := NewPostHandler(nil, mockBiz)

	r := setupGin()
	r.POST("/posts", h.CreatePost)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/posts", nil)
	c.Set("userID", "u1")

	mockPostBiz.On("CreatePost", mock.Anything, mock.AnythingOfType("*model.CreatePostRequest")).Return(nil, errors.New("fail")).Once()

	h.CreatePost(c)
	assert.Equal(t, http.StatusOK, w.Code)
	mockPostBiz.AssertExpectations(t)
}

/**
1、httptest的作用
httptest 是 Go 标准库提供的 HTTP 测试工具，主要用于构造 HTTP 请求（Request）和响应（Response），模拟 HTTP 客户端和服务端的交互。
常用的有 httptest.NewRequest（构造请求）、httptest.NewRecorder（模拟响应写入器）。
2、gin的作用
Gin 是 Web 框架，负责路由分发、参数绑定、上下文管理等。
Gin 的 handler 方法参数是 *gin.Context，它内部封装了 HTTP 请求和响应。
3、两者的关系
方式一： gin路由 + ServeHTTP 全流程测试(集成测试、路由测试)
用 httptest.NewRequest 构造请求，然后用 Gin 的 ServeHTTP 或 CreateTestContext 让 handler 处理这个请求。

方式二：gin.CreateTestConext (handler的单元测试)
httptest 构建的 *http.Request 和 *httptest.ResponseRecorder，可以直接传递给 Gin 的路由或 gin.CreateTestContext，从而驱动 Gin 的 handler 进行测试。
*/

// TestPostHandler_CreatePost_方式一：Gin路由+ServeHTTP全流程测试，（集成测试、路由测试）
func TestPostHandler_CreatePost_FullFlow(t *testing.T) {
	mockPostBiz := new(MockPostBiz)
	mockBiz := &MockBiz{PostBiz: mockPostBiz}
	h := NewPostHandler(nil, mockBiz)

	r := gin.New()
	r.POST("/posts", h.CreatePost)

	reqBody := &model.CreatePostRequest{Title: "test title", Content: "test content"}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/posts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	// Gin中间件通常会设置userID，这里直接用Set方法模拟
	req = req.WithContext(context.WithValue(req.Context(), "userID", "u1"))

	mockPostBiz.On("CreatePost", mock.Anything, mock.AnythingOfType("*model.CreatePostRequest")).Return(&model.Post{PostID: "p1", UserID: "u1"}, nil).Once()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockPostBiz.AssertExpectations(t)
}

// TestPostHandler_CreatePost_方式二：gin.CreateTestContext测试handler(handler单元测试)
func TestPostHandler_CreatePost_ContextOnly(t *testing.T) {
	// mockPostBiz := new(MockPostBiz)
	// mockBiz := &MockBiz{PostBiz: mockPostBiz}
	// h := NewPostHandler(nil, mockBiz)

	postHandlerTest := setupPostHandlerTest()

	reqBody := &model.CreatePostRequest{Title: "test title", Content: "test content"}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/posts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// gin处理HTTP请求，分发到handler
	c.Request = req
	c.Set("userID", "u1")
	postHandlerTest.mockPostBiz.On("CreatePost", mock.Anything, mock.AnythingOfType("*model.CreatePostRequest")).Return(&model.Post{PostID: "p1", UserID: "u1"}, nil).Once()

	postHandlerTest.handler.CreatePost(c)
	assert.Equal(t, http.StatusOK, w.Code)
	postHandlerTest.mockPostBiz.AssertExpectations(t)
}

// 其它接口（DeletePost, GetPostByID, ListPosts, UpdatePost, GetPostsByUserID）可仿照上面写法补充
