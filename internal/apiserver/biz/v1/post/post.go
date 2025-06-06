package biz

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

type PostBiz interface {
	// Create 创建帖子
	CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error)
	// GetByID 根据 ID 获取帖子
	GetPostByID(ctx context.Context, id uint) (*model.Post, error)
	// Update 更新帖子
	UpdatePost(ctx context.Context, post *model.UpdatePostRequest) error
	// Delete 删除帖子
	DeletePost(ctx context.Context, id uint) error
	// List 获取帖子列表
	ListPosts(ctx context.Context, page, pageSize int) (*model.ListPostResponse, error)
	// GetByUserID 根据用户 ID 获取帖子列表
	GetPostsByUserID(ctx context.Context, userID string, page, pageSize int) (*model.ListPostResponse, error)
	// GetByPostID 根据帖子 ID 获取帖子
	GetPostByPostID(ctx context.Context, postID string) (*model.Post, error)
}

// NewPostBiz 实例化postBiz对象
func NewPostBiz(logger *log.Logger, store store.PostStore) PostBiz {

	return &postBiz{
		store: store,
	}
}

// postBiz	实现了post业务层接口
type postBiz struct {
	store store.PostStore
}

// CreatePost implements PostBiz.
func (p *postBiz) CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {
	panic("unimplemented")
}

// DeletePost implements PostBiz.
func (p *postBiz) DeletePost(ctx context.Context, id uint) error {
	panic("unimplemented")
}

// GetPostByID implements PostBiz.
func (p *postBiz) GetPostByID(ctx context.Context, id uint) (*model.Post, error) {
	panic("unimplemented")
}

// GetPostByPostID implements PostBiz.
func (p *postBiz) GetPostByPostID(ctx context.Context, postID string) (*model.Post, error) {
	panic("unimplemented")
}

// GetPostsByUserID implements PostBiz.
func (p *postBiz) GetPostsByUserID(ctx context.Context, userID string, page int, pageSize int) (*model.ListPostResponse, error) {
	panic("unimplemented")
}

// ListPosts implements PostBiz.
func (p *postBiz) ListPosts(ctx context.Context, page int, pageSize int) (*model.ListPostResponse, error) {
	panic("unimplemented")
}

// UpdatePost implements PostBiz.
func (p *postBiz) UpdatePost(ctx context.Context, post *model.UpdatePostRequest) error {
	panic("unimplemented")
}

var _ PostBiz = (*postBiz)(nil)
