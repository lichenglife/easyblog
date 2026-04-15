package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

type CommentStore interface {
	// Create 创建评论
	Create(ctx context.Context, comment *model.Comment) error
	// GetByID 根据 ID 获取评论
	GetByID(ctx context.Context, id uint) (*model.Comment, error)
	// Update 更新评论
	Update(ctx context.Context, comment *model.Comment) error
	// Delete 删除评论
	Delete(ctx context.Context, id uint) error
	// ListByPostID 根据文章 ID 获取评论列表（顶级评论）
	ListByPostID(ctx context.Context, postID uint, page, pageSize int) ([]*model.Comment, error)
	// ListByParentID 根据父评论 ID 获取回复列表（楼中楼）
	ListByParentID(ctx context.Context, parentID uint, page, pageSize int) ([]*model.Comment, error)
	// GetByPostID 根据帖子 ID 获取帖子
	GetByPostID(ctx context.Context, postID uint) (*model.Comment, error)
	// IncrementLikeCount 增加点赞数
	IncrementLikeCount(ctx context.Context, id uint) error
	// IncrementReplyCount 增加回复数
	IncrementReplyCount(ctx context.Context, id uint) error
	// DecrementReplyCount 减少回复数
	DecrementReplyCount(ctx context.Context, id uint) error
	// DecrementLikeCount 减少点赞数
	DecrementLikeCount(ctx context.Context, id uint) error
}

type comments struct {
	db *gorm.DB
}

// NewComments 创建 CommentStore 实例
func NewComments(db *gorm.DB) CommentStore {
	return &comments{db: db}
}

var _ CommentStore = (*comments)(nil)

// Create 创建评论
func (c *comments) Create(ctx context.Context, comment *model.Comment) error {
	return c.db.WithContext(ctx).Create(comment).Error
}

// GetByID 根据 ID 获取评论
func (c *comments) GetByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	err := c.db.WithContext(ctx).First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByPostID 根据帖子 ID 获取帖子
func (c *comments) GetByPostID(ctx context.Context, postID uint) (*model.Comment, error) {
	var comment model.Comment
	err := c.db.WithContext(ctx).Where("post_id = ?", postID).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// Update 更新评论
func (c *comments) Update(ctx context.Context, comment *model.Comment) error {
	return c.db.WithContext(ctx).Save(comment).Error
}

// Delete 删除评论
func (c *comments) Delete(ctx context.Context, id uint) error {
	return c.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

// ListByPostID 根据文章 ID 获取评论列表（顶级评论）
func (c *comments) ListByPostID(ctx context.Context, postID uint, page, pageSize int) ([]*model.Comment, error) {
	var comments []*model.Comment
	offset := (page - 1) * pageSize
	err := c.db.WithContext(ctx).
		Where("post_id = ? AND parent_id = 0 AND status = ?", postID, 1).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// ListByParentID 根据父评论 ID 获取回复列表（楼中楼）
func (c *comments) ListByParentID(ctx context.Context, parentID uint, page, pageSize int) ([]*model.Comment, error) {
	var comments []*model.Comment
	offset := (page - 1) * pageSize
	err := c.db.WithContext(ctx).
		Where("parent_id = ? AND status = ?", parentID, 1).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// IncrementLikeCount 增加点赞数
func (c *comments) IncrementLikeCount(ctx context.Context, id uint) error {
	return c.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", id).
		UpdateColumn("like_count", c.db.Raw("like_count + 1")).Error
}

// IncrementReplyCount 增加回复数
func (c *comments) IncrementReplyCount(ctx context.Context, id uint) error {
	return c.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ?", id).
		UpdateColumn("reply_count", c.db.Raw("reply_count + 1")).Error
}

// DecrementReplyCount 减少回复数
func (c *comments) DecrementReplyCount(ctx context.Context, id uint) error {
	return c.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ? AND reply_count > 0", id).
		UpdateColumn("reply_count", c.db.Raw("CASE WHEN reply_count > 0 THEN reply_count - 1 ELSE 0 END")).Error
}

// DecrementLikeCount 减少点赞数
func (c *comments) DecrementLikeCount(ctx context.Context, id uint) error {
	return c.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("id = ? AND like_count > 0", id).
		UpdateColumn("like_count", c.db.Raw("CASE WHEN like_count > 0 THEN like_count - 1 ELSE 0 END")).Error
}
