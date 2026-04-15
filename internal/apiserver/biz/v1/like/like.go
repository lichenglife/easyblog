package likev1

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"gorm.io/gorm"
)

type LikeBiz interface {
	// LikePost 点赞文章
	LikePost(ctx context.Context, userID, postID uint) error
	// UnlikePost 取消点赞文章
	UnlikePost(ctx context.Context, userID, postID uint) error
	// IsLikedPost 检查文章点赞状态
	IsLikedPost(ctx context.Context, userID, postID uint) (bool, error)
	// GetPostLikeCount 获取文章点赞数
	GetPostLikeCount(ctx context.Context, postID uint) (int64, error)
}

type likeBiz struct {
	likeStore store.LikeStore
	postStore store.PostStore
}

var _ LikeBiz = (*likeBiz)(nil)

func NewLikeBiz(likeStore store.LikeStore, postStore store.PostStore) LikeBiz {
	return &likeBiz{
		likeStore: likeStore,
		postStore: postStore,
	}
}

// LikePost 点赞文章
func (b *likeBiz) LikePost(ctx context.Context, userID, postID uint) error {
	// 检查文章是否存在
	_, err := b.postStore.GetByID(ctx, postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.ErrPostNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 检查是否已点赞
	liked, _ := b.likeStore.IsLiked(ctx, userID, postID, 1) // 1=文章
	if liked {
		return errno.ErrAlreadyLiked
	}

	// 创建点赞记录
	like := &model.Like{
		UserID:     userID,
		TargetID:   postID,
		TargetType: 1, // 1=文章
	}
	if err := b.likeStore.Create(ctx, like); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 更新文章点赞数
	if err := b.postStore.IncrementLikeCount(ctx, postID); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// UnlikePost 取消点赞文章
func (b *likeBiz) UnlikePost(ctx context.Context, userID, postID uint) error {
	// 删除点赞记录
	if err := b.likeStore.Delete(ctx, userID, postID, 1); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 减少文章点赞数
	if err := b.postStore.DecrementLikeCount(ctx, postID); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// IsLikedPost 检查文章点赞状态
func (b *likeBiz) IsLikedPost(ctx context.Context, userID, postID uint) (bool, error) {
	return b.likeStore.IsLiked(ctx, userID, postID, 1)
}

// GetPostLikeCount 获取文章点赞数
func (b *likeBiz) GetPostLikeCount(ctx context.Context, postID uint) (int64, error) {
	return b.likeStore.CountByTarget(ctx, postID, 1)
}
