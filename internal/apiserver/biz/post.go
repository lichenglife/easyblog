package biz

import (
	"context"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/lichenglife/easyblog/internal/pkg/utils"
	"go.uber.org/zap"
)

// postBiz	实现了post业务层接口
type postBiz struct {
	logger *log.Logger
	store  store.PostStore
}

// CreatePost implements PostBiz.
func (p *postBiz) CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {

	post := &model.Post{
		Title:     req.Title,
		Content:   req.Content,
		UserID:    req.UserID,
		PostID:    utils.GenPostID(),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	if err := p.store.Create(ctx, post); err != nil {
		p.logger.Error("create post failed, err: %v", zap.Error(err))
		return nil, err
	}

	return post, nil
}

// DeletePost implements PostBiz.
func (p *postBiz) DeletePostByPostID(ctx context.Context, postID string) error {
	// 根据博客ID 获取博客
	post, err := p.store.GetByPostID(ctx, postID)
	if err != nil {
		p.logger.Error("get post by id failed, err: %v", zap.Error(err))
		return err
	}
	if post == nil {
		p.logger.Logger.Error("post is not  found")
		return errno.ErrNotFound
	}
	// 校验是否与当前用户一致
	userID := ctx.Value("userID").(string)
	if post.UserID != userID {
		p.logger.Logger.Error("post is not belong to the user")
		return errno.ErrPostNotBelongToUser
	}
	// 删除博客
	if err := p.store.Delete(ctx, postID); err != nil {
		p.logger.Error("delete post by id failed, err: %v", zap.Error(err))
		return err
	}
	return nil
}

// GetPostByPostID implements PostBiz.
func (p *postBiz) GetPostByPostID(ctx context.Context, postID string) (*model.Post, error) {
	post, err := p.store.GetByPostID(ctx, postID)
	if err != nil {
		p.logger.Error("get post by id failed, err: %v", zap.Error(err))
		return nil, err
	}
	return post, nil
}

// GetPostsByUserID implements PostBiz.
func (p *postBiz) GetPostsByUserID(ctx context.Context, userID string, page int, pageSize int) (*model.ListPostResponse, error) {
	posts, err := p.store.GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		p.logger.Error("get posts by user id failed, err: %v", zap.Error(err))
		return nil, err
	}
	return &model.ListPostResponse{
		Posts:      posts,
		TotalCount: int64(len(posts)),
		HasMore:    len(posts) == pageSize,
	}, nil
}

// ListPosts implements PostBiz.
func (p *postBiz) ListPosts(ctx context.Context, page int, pageSize int) (*model.ListPostResponse, error) {
	posts, err := p.store.List(ctx, page, pageSize)
	if err != nil {
		p.logger.Error("list posts failed, err: %v", zap.Error(err))
		return nil, err
	}

	return &model.ListPostResponse{
		Posts:      posts,
		TotalCount: int64(len(posts)),
		HasMore:    len(posts) == pageSize,
	}, nil
}

// UpdatePost implements PostBiz.
func (p *postBiz) UpdatePost(ctx context.Context, updatePost *model.UpdatePostRequest) error {
	// 查询判断是否存在
	post, err := p.store.GetByPostID(ctx, updatePost.PostID)
	if err != nil {
		p.logger.Error("get post by id failed, err: %v", zap.Error(err))
		return err
	}
	if post == nil {
		p.logger.Logger.Error("post is not found")
		return errno.ErrNotFound
	}
	// 更新博客
	post.Title = updatePost.Title
	post.Content = updatePost.Content
	post.UpdatedAt = time.Now()
	// 更新博客
	if err := p.store.Update(ctx, post); err != nil {
		p.logger.Error("update post failed, err: %v", zap.Error(err))
		return err
	}
	return nil
}

var _ PostBiz = (*postBiz)(nil)
