package biz

import (
	"context"
	"errors"
	"testing"

	// "time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func TestMain(m *testing.M) {
// 	// 用 zap 的 NewNop() 创建一个不会输出的 logger
// 	log.Log = &log.Logger{zap.NewNop()}
// 	m.Run()
// }

// MockPostStore 是对 store.PostStore 的 mock 实现，用于单元测试 biz 层时模拟数据层行为。
// 只需实现用到的方法即可。
type MockPostStore struct {
	mock.Mock
}

// MockPostStore 实现store.PostStore
func (m *MockPostStore) Create(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}
func (m *MockPostStore) Delete(ctx context.Context, postID string) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}
func (m *MockPostStore) GetByPostID(ctx context.Context, postID string) (*model.Post, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}
func (m *MockPostStore) GetByUserID(ctx context.Context, userID string, page, pageSize int) (int64, []*model.Post, error) {
	args := m.Called(ctx, userID, page, pageSize)
	return args.Get(0).(int64), args.Get(1).([]*model.Post), args.Error(2)
}
func (m *MockPostStore) List(ctx context.Context, page, pageSize int) (int64, []*model.Post, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).(int64), args.Get(1).([]*model.Post), args.Error(2)
}
func (m *MockPostStore) Update(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}
func (m *MockPostStore) GetByID(ctx context.Context, id uint) (*model.Post, error) { return nil, nil }

// MockStore 是对 store.IStore 的 mock 实现，只需实现 Post() 方法。
type MockStore struct {
	mock.Mock
	PostStore store.PostStore
}

func (m *MockStore) Post() store.PostStore {
	return m.PostStore
}
func (m *MockStore) User() store.UserStore { return nil }
func (m *MockStore) Close() error          { return nil }

// TestPostBiz_CreatePost 测试创建博客的正常流程。
func TestPostBiz_CreatePost(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	req := &model.CreatePostRequest{
		Title:   "test",
		Content: "content",
		UserID:  "u1",
	}
	// 期望 Create 被调用一次，返回 nil
	mockPostStore.On("Create", mock.Anything, mock.AnythingOfType("*model.Post")).Return(nil).Once()

	post, err := biz.CreatePost(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, req.Title, post.Title)
	assert.Equal(t, req.Content, post.Content)
	assert.Equal(t, req.UserID, post.UserID)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_CreatePost_Fail 测试创建博客失败场景。
func TestPostBiz_CreatePost_Fail(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	req := &model.CreatePostRequest{
		Title:   "test",
		Content: "content",
		UserID:  "u1",
	}
	// 期望 Create 返回错误
	mockPostStore.On("Create", mock.Anything, mock.AnythingOfType("*model.Post")).Return(errors.New("db error")).Once()

	post, err := biz.CreatePost(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, post)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_DeletePostByPostID 测试删除博客的正常流程。
func TestPostBiz_DeletePostByPostID(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	ctx := context.WithValue(context.Background(), "userID", "u1")
	post := &model.Post{PostID: "p1", UserID: "u1"}
	mockPostStore.On("GetByPostID", mock.Anything, "p1").Return(post, nil).Once()
	mockPostStore.On("Delete", mock.Anything, "p1").Return(nil).Once()

	err := biz.DeletePostByPostID(ctx, "p1")
	assert.NoError(t, err)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_DeletePostByPostID_NotFound 测试删除博客时未找到博客。
func TestPostBiz_DeletePostByPostID_NotFound(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	ctx := context.WithValue(context.Background(), "userID", "u1")
	mockPostStore.On("GetByPostID", mock.Anything, "p1").Return(nil, errors.New("not found")).Once()

	err := biz.DeletePostByPostID(ctx, "p1")
	assert.Error(t, err)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_DeletePostByPostID_NotBelongToUser 测试删除博客时博客不属于当前用户。
func TestPostBiz_DeletePostByPostID_NotBelongToUser(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	ctx := context.WithValue(context.Background(), "userID", "u1")
	post := &model.Post{PostID: "p1", UserID: "u2"}
	mockPostStore.On("GetByPostID", mock.Anything, "p1").Return(post, nil).Once()

	err := biz.DeletePostByPostID(ctx, "p1")
	assert.Error(t, err)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_GetPostByPostID 测试根据 PostID 查询博客。
func TestPostBiz_GetPostByPostID(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	post := &model.Post{PostID: "p1", UserID: "u1"}
	mockPostStore.On("GetByPostID", mock.Anything, "p1").Return(post, nil).Once()

	got, err := biz.GetPostByPostID(context.Background(), "p1")
	assert.NoError(t, err)
	assert.Equal(t, post, got)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_GetPostsByUserID 测试查询用户所有博客。
func TestPostBiz_GetPostsByUserID(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	posts := []*model.Post{{PostID: "p1", UserID: "u1"}}
	mockPostStore.On("GetByUserID", mock.Anything, "u1", 1, 10).Return(int64(1), posts, nil).Once()

	resp, err := biz.GetPostsByUserID(context.Background(), "u1", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.TotalCount)
	assert.Equal(t, posts, resp.Posts)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_ListPosts 测试分页查询所有博客。
func TestPostBiz_ListPosts(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	posts := []*model.Post{{PostID: "p1", UserID: "u1"}}
	mockPostStore.On("List", mock.Anything, 1, 10).Return(int64(1), posts, nil).Once()

	resp, err := biz.ListPosts(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.TotalCount)
	assert.Equal(t, posts, resp.Posts)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_UpdatePost 测试更新博客的正常流程。
func TestPostBiz_UpdatePost(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	post := &model.Post{PostID: "p1", UserID: "u1", Title: "old", Content: "old"}
	updateReq := &model.UpdatePostRequest{PostID: "p1", Title: "new", Content: "new"}
	mockPostStore.On("GetByPostID", mock.Anything, "p1").Return(post, nil).Once()
	mockPostStore.On("Update", mock.Anything, mock.AnythingOfType("*model.Post")).Return(nil).Once()

	err := biz.UpdatePost(context.Background(), updateReq)
	assert.NoError(t, err)
	mockPostStore.AssertExpectations(t)
}

// TestPostBiz_UpdatePost_NotFound 测试更新博客时未找到博客。
func TestPostBiz_UpdatePost_NotFound(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockStore := &MockStore{PostStore: mockPostStore}
	biz := NewPostBiz(mockStore)

	updateReq := &model.UpdatePostRequest{PostID: "p1", Title: "new", Content: "new"}
	mockPostStore.On("GetByPostID", mock.Anything, "p1").Return(nil, errors.New("not found")).Once()

	err := biz.UpdatePost(context.Background(), updateReq)
	assert.Error(t, err)
	mockPostStore.AssertExpectations(t)
}

/**
技术说明
1. biz 层依赖 store 层的 mock 技术
- 接口注入：biz 层依赖 store 层接口（如 store.IStore），测试时注入 mock 实现。
- testify/mock：通过实现 mock 结构体，模拟 store 层方法的返回值和行为。
mock.On().Return()：设定期望的调用参数和返回值，灵活模拟各种分支。
2. mock store 及 mock 返回值实现
- MockPostStore 实现了 store.PostStore，每个方法都用 m.Called(...) 记录调用和返回。
- MockStore 实现了 store.IStore，只需实现 Post() 方法返回 MockPostStore。
- 在测试用例中通过 mockPostStore.On(...).Return(...) 设定返回值。
3. 断言
-使用 assert 断言 biz 层的返回值和错误，确保逻辑正确。
- 使用 mockPostStore.AssertExpectations(t) 保证 mock 方法被正确调用。

*/
