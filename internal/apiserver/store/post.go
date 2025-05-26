package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

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
func (p *posts) GetByID(ctx context.Context, postID string) (*model.Post, error) {
	return nil, nil
}

// Update 更新帖子
func (p *posts) Update(ctx context.Context, post *model.Post) error {
	return p.db.Model(&model.Post{}).Where("postID = ?", post.PostID).Updates(post).Error
}

// Delete 删除帖子
func (p *posts) Delete(ctx context.Context, postID string) error {
	err := p.db.Delete(&model.Post{}, "postID = ?", postID).Error
	if err != nil {
		return err
	}
	return nil
}

// List 获取帖子列表
func (p *posts) List(ctx context.Context, page, pageSize int) ([]*model.Post, error) {
	var count int64
	if err := p.db.Model(&model.Post{}).Count(&count).Error; err != nil {
		return nil, err
	}
	var posts []*model.Post

	if err := p.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// GetByUserID 根据用户 ID 获取帖子列表
func (p *posts) GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.Post, error) {
	posts := make([]*model.Post, 0)
	err := p.db.Where("userID = ?", userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// GetByPostID 根据帖子 ID 获取帖子
func (p *posts) GetByPostID(ctx context.Context, postID string) (*model.Post, error) {
	post := &model.Post{}
	err := p.db.Where("postID = ?", postID).First(post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}
