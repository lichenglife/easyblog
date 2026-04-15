package biz

import (
	"context"
	"testing"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostStore 模拟 PostStore
type MockPostStore struct {
	mock.Mock
}

func (m *MockPostStore) Create(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostStore) GetByID(ctx context.Context, id uint) (*model.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostStore) GetByPostID(ctx context.Context, postID string) (*model.Post, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostStore) Update(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostStore) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostStore) List(ctx context.Context, page int, pageSize int) ([]*model.Post, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostStore) ListByUserID(ctx context.Context, userID uint, page int, pageSize int) ([]*model.Post, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostStore) GetCategories(ctx context.Context, parentID *uint) ([]*model.Category, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Category), args.Error(1)
}

// SearchPosts 搜索文章
func (m *MockPostStore) SearchPosts(ctx context.Context, keyword string, page, pageSize int) ([]*model.Post, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

// ListByCategoryID 根据分类 ID 获取文章列表
func (m *MockPostStore) ListByCategoryID(ctx context.Context, categoryID uint, page, pageSize int) ([]*model.Post, error) {
	args := m.Called(ctx, categoryID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

// ListByTagID 根据标签 ID 获取文章列表
func (m *MockPostStore) ListByTagID(ctx context.Context, tagID uint, page, pageSize int) ([]*model.Post, error) {
	args := m.Called(ctx, tagID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

// SetTop 置顶文章
func (m *MockPostStore) SetTop(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CancelTop 取消置顶文章
func (m *MockPostStore) CancelTop(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// IncrementViewCount 增加阅读量
func (m *MockPostStore) IncrementViewCount(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// IncrementLikeCount 增加点赞数
func (m *MockPostStore) IncrementLikeCount(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DecrementLikeCount 减少点赞数
func (m *MockPostStore) DecrementLikeCount(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockPostTagStore 模拟 PostTagStore
type MockPostTagStore struct {
	mock.Mock
}

func (m *MockPostTagStore) Create(ctx context.Context, postTag *model.PostTag) error {
	args := m.Called(ctx, postTag)
	return args.Error(0)
}

func (m *MockPostTagStore) CreateBatch(ctx context.Context, postTags []*model.PostTag) error {
	args := m.Called(ctx, postTags)
	return args.Error(0)
}

func (m *MockPostTagStore) DeleteByPostID(ctx context.Context, postID uint) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockPostTagStore) DeleteByTagID(ctx context.Context, tagID uint) error {
	args := m.Called(ctx, tagID)
	return args.Error(0)
}

func (m *MockPostTagStore) ListByPostID(ctx context.Context, postID uint) ([]model.PostTag, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PostTag), args.Error(1)
}

func (m *MockPostTagStore) ListTagIDsByPostID(ctx context.Context, postID uint) ([]uint, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]uint), args.Error(1)
}

func (m *MockPostTagStore) ListByTagID(ctx context.Context, tagID uint) ([]model.PostTag, error) {
	args := m.Called(ctx, tagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PostTag), args.Error(1)
}

func (m *MockPostTagStore) ListPostIDsByTagID(ctx context.Context, tagID uint, page, pageSize int) ([]uint, int64, error) {
	args := m.Called(ctx, tagID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]uint), args.Get(1).(int64), args.Error(2)
}

func newTestPostBiz(mockStore *MockPostStore, mockPostTagStore *MockPostTagStore) PostBiz {
	return NewPostBiz(mockStore, mockPostTagStore)
}

func TestPostBiz_CreatePost_Success(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()
	req := &model.CreatePostRequest{
		Title:   "Test Post",
		Content: "Test content",
		Summary: "Test summary",
		TagIDs:  []uint{1, 2, 3}, // 添加标签 ID
	}

	mockStore.On("Create", ctx, mock.AnythingOfType("*model.Post")).Return(nil).Run(func(args mock.Arguments) {
		// Verify the post passed to Create
		post := args.Get(1).(*model.Post)
		assert.Equal(t, uint(1), post.UserID)
		assert.Equal(t, req.Title, post.Title)
		assert.Equal(t, req.Content, post.Content)
	})

	mockPostTagStore.On("CreateBatch", ctx, mock.AnythingOfType("[]*model.PostTag")).Return(nil)

	result, err := biz.CreatePost(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Content, result.Content)
	assert.Equal(t, req.TagIDs, result.Tags)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_CreatePost_DatabaseError(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()
	req := &model.CreatePostRequest{
		Title:   "Test Post",
		Content: "Test content",
	}

	mockStore.On("Create", ctx, mock.AnythingOfType("*model.Post")).Return(errno.ErrDatabase)

	result, err := biz.CreatePost(ctx, 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NotNil(t, err)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_GetPostByID_Success(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()
	mockPost := &model.Post{
		ID:      1,
		PostID:  "test-uuid",
		UserID:  1,
		Title:   "Test Post",
		Content: "Test content",
		Summary: "Test summary",
	}

	mockStore.On("GetByID", ctx, uint(1)).Return(mockPost, nil)
	mockPostTagStore.On("ListTagIDsByPostID", ctx, uint(1)).Return([]uint{1, 2, 3}, nil)

	result, err := biz.GetPostByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockPost.Title, result.Title)
	assert.Equal(t, mockPost.Content, result.Content)
	assert.Equal(t, mockPost.PostID, result.PostID)
	assert.Equal(t, []uint{1, 2, 3}, result.Tags)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_GetPostByID_NotFound(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()

	mockStore.On("GetByID", ctx, uint(999)).Return((*model.Post)(nil), errno.ErrPostNotFound)

	result, err := biz.GetPostByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_GetPostByPostID_Success(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()
	mockPost := &model.Post{
		ID:      1,
		PostID:  "test-uuid",
		UserID:  1,
		Title:   "Test Post",
		Content: "Test content",
	}

	mockStore.On("GetByPostID", ctx, "test-uuid").Return(mockPost, nil)

	result, err := biz.GetPostByPostID(ctx, "test-uuid")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockPost.PostID, result.PostID)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_GetPostByPostID_NotFound(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()

	mockStore.On("GetByPostID", ctx, "non-existent").Return((*model.Post)(nil), errno.ErrPostNotFound)

	result, err := biz.GetPostByPostID(ctx, "non-existent")

	assert.Error(t, err)
	assert.Nil(t, result)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_UpdatePost_Success(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()
	existingPost := &model.Post{
		ID:      1,
		PostID:  "test-uuid",
		UserID:  1,
		Title:   "Old Title",
		Content: "Old content",
	}

	req := &model.UpdatePostRequest{
		ID:      1,
		Title:   "New Title",
		Content: "New content",
		Summary: "New summary",
		TagIDs:  []uint{1, 2}, // 添加标签 ID
	}

	mockStore.On("GetByID", ctx, uint(1)).Return(existingPost, nil)
	mockStore.On("Update", ctx, mock.AnythingOfType("*model.Post")).Return(nil).Run(func(args mock.Arguments) {
		post := args.Get(1).(*model.Post)
		assert.Equal(t, req.Title, post.Title)
		assert.Equal(t, req.Content, post.Content)
		assert.Equal(t, req.Summary, post.Summary)
	})
	mockPostTagStore.On("DeleteByPostID", ctx, uint(1)).Return(nil)
	mockPostTagStore.On("CreateBatch", ctx, mock.AnythingOfType("[]*model.PostTag")).Return(nil)

	err := biz.UpdatePost(ctx, 1, req)

	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_UpdatePost_NotFound(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()
	req := &model.UpdatePostRequest{
		ID:    999,
		Title: "New Title",
	}

	mockStore.On("GetByID", ctx, uint(999)).Return((*model.Post)(nil), errno.ErrPostNotFound)

	err := biz.UpdatePost(ctx, 999, req)

	assert.Error(t, err)

	mockStore.AssertExpectations(t)
}

func TestPostBiz_DeletePost_Success(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()

	mockStore.On("GetByID", ctx, uint(1)).Return(&model.Post{ID: 1}, nil)
	mockPostTagStore.On("DeleteByPostID", ctx, uint(1)).Return(nil)
	mockStore.On("Delete", ctx, uint(1)).Return(nil)

	err := biz.DeletePost(ctx, 1)

	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_DeletePost_Error(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()

	mockStore.On("GetByID", ctx, uint(999)).Return(&model.Post{ID: 999}, nil)
	mockPostTagStore.On("DeleteByPostID", ctx, uint(999)).Return(nil)
	mockStore.On("Delete", ctx, uint(999)).Return(errno.ErrDatabase)

	err := biz.DeletePost(ctx, 999)

	assert.Error(t, err)
	assert.NotNil(t, err)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_ListPosts_Success(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()
	mockPosts := []*model.Post{
		{ID: 1, PostID: "uuid-1", Title: "Post 1", Content: "Content 1"},
		{ID: 2, PostID: "uuid-2", Title: "Post 2", Content: "Content 2"},
		{ID: 3, PostID: "uuid-3", Title: "Post 3", Content: "Content 3"},
	}

	mockStore.On("List", ctx, 1, 10).Return(mockPosts, nil)
	mockPostTagStore.On("ListTagIDsByPostID", ctx, uint(1)).Return([]uint{}, nil)
	mockPostTagStore.On("ListTagIDsByPostID", ctx, uint(2)).Return([]uint{}, nil)
	mockPostTagStore.On("ListTagIDsByPostID", ctx, uint(3)).Return([]uint{}, nil)

	result, err := biz.ListPosts(ctx, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Posts, 3)
	assert.Equal(t, "Post 1", result.Posts[0].Title)
	assert.Equal(t, "Post 2", result.Posts[1].Title)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_ListPosts_Empty(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()

	mockStore.On("List", ctx, 1, 10).Return([]*model.Post{}, nil)

	result, err := biz.ListPosts(ctx, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Posts)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_GetPostsByUserID_Success(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()
	mockPosts := []*model.Post{
		{ID: 1, PostID: "uuid-1", UserID: 1, Title: "My Post 1", Content: "Content 1"},
		{ID: 2, PostID: "uuid-2", UserID: 1, Title: "My Post 2", Content: "Content 2"},
	}

	mockStore.On("ListByUserID", ctx, uint(1), 1, 10).Return(mockPosts, nil)
	mockPostTagStore.On("ListTagIDsByPostID", ctx, uint(1)).Return([]uint{}, nil)
	mockPostTagStore.On("ListTagIDsByPostID", ctx, uint(2)).Return([]uint{}, nil)

	result, err := biz.GetPostsByUserID(ctx, 1, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Posts, 2)
	assert.Equal(t, uint(1), result.Posts[0].UserID)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}

func TestPostBiz_GetPostsByUserID_Empty(t *testing.T) {
	mockStore := new(MockPostStore)
	mockPostTagStore := new(MockPostTagStore)
	biz := newTestPostBiz(mockStore, mockPostTagStore)

	ctx := context.Background()

	mockStore.On("ListByUserID", ctx, uint(1), 1, 10).Return([]*model.Post{}, nil)

	result, err := biz.GetPostsByUserID(ctx, 1, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Posts)

	mockStore.AssertExpectations(t)
	mockPostTagStore.AssertExpectations(t)
}
