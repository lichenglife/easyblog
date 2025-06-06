package biz

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// postBiz	实现了post业务层接口
type postBiz struct {
	logger *log.Logger
	store  store.PostStore
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
