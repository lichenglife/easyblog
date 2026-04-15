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
	// ListByUserID 根据用户 ID 获取帖子列表
	ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*model.Post, error)
	// GetByPostID 根据帖子 ID 获取帖子
	GetByPostID(ctx context.Context, postID string) (*model.Post, error)
	// SearchPosts 搜索文章
	SearchPosts(ctx context.Context, keyword string, page, pageSize int) ([]*model.Post, error)
	// ListByCategoryID 根据分类 ID 获取文章列表
	ListByCategoryID(ctx context.Context, categoryID uint, page, pageSize int) ([]*model.Post, error)
	// ListByTagID 根据标签 ID 获取文章列表
	ListByTagID(ctx context.Context, tagID uint, page, pageSize int) ([]*model.Post, error)
	// SetTop 置顶文章
	SetTop(ctx context.Context, id uint) error
	// CancelTop 取消置顶文章
	CancelTop(ctx context.Context, id uint) error
	// IncrementViewCount 增加阅读量
	IncrementViewCount(ctx context.Context, id uint) error
	// IncrementLikeCount 增加点赞数
	IncrementLikeCount(ctx context.Context, id uint) error
	// DecrementLikeCount 减少点赞数
	DecrementLikeCount(ctx context.Context, id uint) error
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
	return p.db.WithContext(ctx).Create(post).Error
}

// GetByID 根据 ID 获取帖子
func (p *posts) GetByID(ctx context.Context, id uint) (*model.Post, error) {
	var post model.Post
	err := p.db.WithContext(ctx).First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Update 更新帖子
func (p *posts) Update(ctx context.Context, post *model.Post) error {
	return p.db.WithContext(ctx).Save(post).Error
}

// Delete 删除帖子
func (p *posts) Delete(ctx context.Context, id uint) error {
	return p.db.WithContext(ctx).Delete(&model.Post{}, id).Error
}

// List 获取帖子列表
func (p *posts) List(ctx context.Context, page, pageSize int) ([]*model.Post, error) {
	var posts []*model.Post
	offset := (page - 1) * pageSize
	err := p.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// ListByUserID 根据用户 ID 获取帖子列表
func (p *posts) ListByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*model.Post, error) {
	var posts []*model.Post
	offset := (page - 1) * pageSize
	err := p.db.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// GetByPostID 根据帖子 ID 获取帖子
func (p *posts) GetByPostID(ctx context.Context, postID string) (*model.Post, error) {
	var post model.Post
	err := p.db.WithContext(ctx).Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// SearchPosts 搜索文章 - 根据标题和内容搜索
func (p *posts) SearchPosts(ctx context.Context, keyword string, page, pageSize int) ([]*model.Post, error) {
	var posts []*model.Post
	offset := (page - 1) * pageSize
	err := p.db.WithContext(ctx).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Offset(offset).Limit(pageSize).
		Order("is_top DESC, created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// ListByCategoryID 根据分类 ID 获取文章列表
func (p *posts) ListByCategoryID(ctx context.Context, categoryID uint, page, pageSize int) ([]*model.Post, error) {
	var posts []*model.Post
	offset := (page - 1) * pageSize
	err := p.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Offset(offset).Limit(pageSize).
		Order("is_top DESC, created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// ListByTagID 根据标签 ID 获取文章列表 - 需要通过 post_tags 关联表
func (p *posts) ListByTagID(ctx context.Context, tagID uint, page, pageSize int) ([]*model.Post, error) {
	var posts []*model.Post
	offset := (page - 1) * pageSize
	err := p.db.WithContext(ctx).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagID).
		Offset(offset).Limit(pageSize).
		Order("post_tags.created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// SetTop 置顶文章
func (p *posts) SetTop(ctx context.Context, id uint) error {
	return p.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).Update("is_top", 1).Error
}

// CancelTop 取消置顶文章
func (p *posts) CancelTop(ctx context.Context, id uint) error {
	return p.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).Update("is_top", 0).Error
}

// IncrementViewCount 增加阅读量
func (p *posts) IncrementViewCount(ctx context.Context, id uint) error {
	return p.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		UpdateColumn("view_count", p.db.Raw("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞数
func (p *posts) IncrementLikeCount(ctx context.Context, id uint) error {
	return p.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ?", id).
		UpdateColumn("like_count", p.db.Raw("like_count + 1")).Error
}

// DecrementLikeCount 减少点赞数
func (p *posts) DecrementLikeCount(ctx context.Context, id uint) error {
	return p.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("id = ? AND like_count > 0", id).
		UpdateColumn("like_count", p.db.Raw("CASE WHEN like_count > 0 THEN like_count - 1 ELSE 0 END")).Error
}
