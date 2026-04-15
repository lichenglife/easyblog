package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

type LikeStore interface {
	// Create 创建点赞
	Create(ctx context.Context, like *model.Like) error
	// Delete 取消点赞
	Delete(ctx context.Context, userID, targetID, targetType uint) error
	// Get 获取点赞记录
	Get(ctx context.Context, userID, targetID, targetType uint) (*model.Like, error)
	// IsLiked 检查是否已点赞
	IsLiked(ctx context.Context, userID, targetID, targetType uint) (bool, error)
	// CountByTarget 统计目标的点赞数
	CountByTarget(ctx context.Context, targetID, targetType uint) (int64, error)
}

type likes struct {
	db *gorm.DB
}

// NewLikes 创建 LikeStore 实例
func NewLikes(db *gorm.DB) LikeStore {
	return &likes{db: db}
}

var _ LikeStore = (*likes)(nil)

// Create 创建点赞
func (l *likes) Create(ctx context.Context, like *model.Like) error {
	return l.db.WithContext(ctx).Create(like).Error
}

// Delete 取消点赞
func (l *likes) Delete(ctx context.Context, userID, targetID, targetType uint) error {
	return l.db.WithContext(ctx).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Delete(&model.Like{}).Error
}

// Get 获取点赞记录
func (l *likes) Get(ctx context.Context, userID, targetID, targetType uint) (*model.Like, error) {
	var like model.Like
	err := l.db.WithContext(ctx).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// IsLiked 检查是否已点赞
func (l *likes) IsLiked(ctx context.Context, userID, targetID, targetType uint) (bool, error) {
	var count int64
	err := l.db.WithContext(ctx).
		Model(&model.Like{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByTarget 统计目标的点赞数
func (l *likes) CountByTarget(ctx context.Context, targetID, targetType uint) (int64, error) {
	var count int64
	err := l.db.WithContext(ctx).
		Model(&model.Like{}).
		Where("target_id = ? AND target_type = ?", targetID, targetType).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
