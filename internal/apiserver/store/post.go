package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

type PostStore interface {
	// Create 创建帖子
	Create(ctx context.Context, post *model.Post) error
	// GetByID 根据 ID 获取帖子
	GetByID(ctx context.Context, id uint) (*model.Post, error)
	// Update 更新帖子
	Update(ctx context.Context, post *model.Post) error
	// Delete 删除帖子
	Delete(ctx context.Context, id uint) error
	// List 获取帖子列表
	List(ctx context.Context, page, pageSize int) ([]*model.Post, error)
	// GetByUserID 根据用户 ID 获取帖子列表
	GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.Post, error)
	// GetByPostID 根据帖子 ID 获取帖子
	GetByPostID(ctx context.Context, postID string) (*model.Post, error)
}

// postStore 实现Factory 的全部接口
type posts struct {
	db *gorm.DB
}

// newPostStore 创建 postStore 实例
func NewPosts(db *gorm.DB) PostStore {
	return &posts{db: db}
}

// Create 创建帖子
func (p *posts) Create(ctx context.Context, post *model.Post) error {
	return nil
}

// GetByID 根据 ID 获取帖子
func (p *posts) GetByID(ctx context.Context, id uint) (*model.Post, error) {
	return nil, nil
}

// Update 更新帖子
func (p *posts) Update(ctx context.Context, post *model.Post) error {
	return nil
}

// Delete 删除帖子
func (p *posts) Delete(ctx context.Context, id uint) error {
	return nil
}

// List 获取帖子列表
func (p *posts) List(ctx context.Context, page, pageSize int) ([]*model.Post, error) {
	return nil, nil
}

// GetByUserID 根据用户 ID 获取帖子列表
func (p *posts) GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.Post, error) {
	return nil, nil
}

// GetByPostID 根据帖子 ID 获取帖子
func (p *posts) GetByPostID(ctx context.Context, postID string) (*model.Post, error) {
	return nil, nil
}
