package biz

import (
	"context"
	"errors"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	genid "github.com/lichenglife/easyblog/internal/pkg/utils/genID"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostBiz interface {
	// 创建博客
	CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error)
	// 删除博客
	DeletePostByPostID(ctx context.Context, postID string) error
	// 查询博客
	GetPostByPostID(ctx context.Context, postID string) (*model.Post, error)
	// 查询用户名下所有博客
	GetPostsByUserID(ctx context.Context, userID string, page int, pageSize int) (*model.ListPostResponse, error)
	// 查询博客
	ListPosts(ctx context.Context, page int, pageSize int) (*model.ListPostResponse, error)
	// 更新
	UpdatePost(ctx context.Context, updatePost *model.UpdatePostRequest) error
}

// postBiz	实现了post业务层接口
type postBiz struct {
	store store.IStore
}

func NewPostBiz(store store.IStore) PostBiz {
	return &postBiz{
		store: store,
	}
}

// CreatePost implements PostBiz.
func (p *postBiz) CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {

	post := &model.Post{
		Title:     req.Title,
		Content:   req.Content,
		UserID:    req.UserID,
		PostID:    genid.GeneratePostID(),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	if err := p.store.Post().Create(ctx, post); err != nil {
		log.Log.Error("create post failed, err: %v", zap.Error(err))
		return nil, err
	}

	return post, nil
}

// DeletePost implements PostBiz.
func (p *postBiz) DeletePostByPostID(ctx context.Context, postID string) error {
	// 根据博客ID 获取博客
	post, err := p.store.Post().GetByPostID(ctx, postID)
	if err != nil {
		log.Log.Error("get post by id failed, err: %v", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrNotFound
		}
	}
	if post == nil {
		log.Log.Logger.Error("post is not  found")
		return errno.ErrNotFound
	}
	// 校验是否与当前用户一致
	userID := ctx.Value("userID").(string)
	if post.UserID != userID {
		log.Log.Logger.Error("post is not belong to the user")
		return errno.ErrPostNotBelongToUser
	}
	// 删除博客
	if err := p.store.Post().Delete(ctx, postID); err != nil {
		log.Log.Error("delete post by id failed, err: %v", zap.Error(err))
		return err
	}
	return nil
}

// GetPostByPostID implements PostBiz.
func (p *postBiz) GetPostByPostID(ctx context.Context, postID string) (*model.Post, error) {
	post, err := p.store.Post().GetByPostID(ctx, postID)
	if err != nil {
		log.Log.Error("get post by id failed, err: %v", zap.Error(err))
		return nil, err
	}
	return post, nil
}

// GetPostsByUserID implements PostBiz.
func (p *postBiz) GetPostsByUserID(ctx context.Context, userID string, page int, pageSize int) (*model.ListPostResponse, error) {

	count, posts, err := p.store.Post().GetByUserID(ctx, userID, page, pageSize)
	if err != nil {
		log.Log.Error("get posts by user id failed, err: %v", zap.Error(err))
		return nil, err
	}
	return &model.ListPostResponse{
		Posts:      posts,
		TotalCount: int64(len(posts)),
		HasMore:    int(count) > len(posts),
	}, nil
}

// ListPosts implements PostBiz.
func (p *postBiz) ListPosts(ctx context.Context, page int, pageSize int) (*model.ListPostResponse, error) {
	count, posts, err := p.store.Post().List(ctx, page, pageSize)
	if err != nil {
		log.Log.Error("list posts failed, err: %v", zap.Error(err))
		return nil, err
	}

	return &model.ListPostResponse{
		Posts:      posts,
		TotalCount: int64(len(posts)),
		HasMore:    int(count) > pageSize,
	}, nil
}

// UpdatePost implements PostBiz.
func (p *postBiz) UpdatePost(ctx context.Context, updatePost *model.UpdatePostRequest) error {
	// 查询判断是否存在
	post, err := p.store.Post().GetByPostID(ctx, updatePost.PostID)
	if err != nil {
		log.Log.Error("failed to get post %v", zap.Error(err))
		return err
	}
	if post == nil {
		return errno.ErrNotFound
	}
	// 更新博客
	post.Title = updatePost.Title
	post.Content = updatePost.Content
	post.UpdatedAt = time.Now()
	// 更新博客
	if err := p.store.Post().Update(ctx, post); err != nil {
		log.Log.Error("update post failed, err: %v", zap.Error(err))
		return err
	}
	return nil
}

var _ PostBiz = (*postBiz)(nil)
