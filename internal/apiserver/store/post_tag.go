package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

type PostTagStore interface {
	// Create 创建文章 - 标签关联
	Create(ctx context.Context, postTag *model.PostTag) error
	// CreateBatch 批量创建文章 - 标签关联
	CreateBatch(ctx context.Context, postTags []*model.PostTag) error
	// DeleteByPostID 根据文章 ID 删除关联
	DeleteByPostID(ctx context.Context, postID uint) error
	// DeleteByTagID 根据标签 ID 删除关联
	DeleteByTagID(ctx context.Context, tagID uint) error
	// ListByPostID 根据文章 ID 获取标签列表
	ListByPostID(ctx context.Context, postID uint) ([]model.PostTag, error)
	// ListTagIDsByPostID 根据文章 ID 获取标签 ID 列表
	ListTagIDsByPostID(ctx context.Context, postID uint) ([]uint, error)
	// ListByTagID 根据标签 ID 获取文章关联列表
	ListByTagID(ctx context.Context, tagID uint) ([]model.PostTag, error)
	// ListPostIDsByTagID 根据标签 ID 获取文章 ID 列表（分页）
	ListPostIDsByTagID(ctx context.Context, tagID uint, page, pageSize int) ([]uint, int64, error)
}

// postTagStore 实现 PostTagStore 接口
type postTagStore struct {
	db           *gorm.DB
	postTagStore *posts
}

// NewPostTagStore 创建 PostTagStore 实例
func NewPostTagStore(db *gorm.DB) PostTagStore {
	return &postTagStore{
		db:           db,
		postTagStore: &posts{db: db},
	}
}

// Create 创建文章 - 标签关联
func (p *postTagStore) Create(ctx context.Context, postTag *model.PostTag) error {
	return p.db.WithContext(ctx).Create(postTag).Error
}

// CreateBatch 批量创建文章 - 标签关联
func (p *postTagStore) CreateBatch(ctx context.Context, postTags []*model.PostTag) error {
	return p.db.WithContext(ctx).CreateInBatches(postTags, len(postTags)).Error
}

// DeleteByPostID 根据文章 ID 删除关联
func (p *postTagStore) DeleteByPostID(ctx context.Context, postID uint) error {
	return p.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&model.PostTag{}).Error
}

// DeleteByTagID 根据标签 ID 删除关联
func (p *postTagStore) DeleteByTagID(ctx context.Context, tagID uint) error {
	return p.db.WithContext(ctx).Where("tag_id = ?", tagID).Delete(&model.PostTag{}).Error
}

// ListByPostID 根据文章 ID 获取标签列表
func (p *postTagStore) ListByPostID(ctx context.Context, postID uint) ([]model.PostTag, error) {
	var postTags []model.PostTag
	err := p.db.WithContext(ctx).Where("post_id = ?", postID).Find(&postTags).Error
	if err != nil {
		return nil, err
	}
	return postTags, nil
}

// ListTagIDsByPostID 根据文章 ID 获取标签 ID 列表
func (p *postTagStore) ListTagIDsByPostID(ctx context.Context, postID uint) ([]uint, error) {
	var tagIDs []uint
	err := p.db.WithContext(ctx).
		Model(&model.PostTag{}).
		Where("post_id = ?", postID).
		Pluck("tag_id", &tagIDs).Error
	return tagIDs, err
}

// ListByTagID 根据标签 ID 获取文章关联列表
func (p *postTagStore) ListByTagID(ctx context.Context, tagID uint) ([]model.PostTag, error) {
	var postTags []model.PostTag
	err := p.db.WithContext(ctx).Where("tag_id = ?", tagID).Find(&postTags).Error
	if err != nil {
		return nil, err
	}
	return postTags, nil
}

// ListPostIDsByTagID 根据标签 ID 获取文章 ID 列表（分页）
func (p *postTagStore) ListPostIDsByTagID(ctx context.Context, tagID uint, page, pageSize int) ([]uint, int64, error) {
	var postIDs []uint
	var total int64
	offset := (page - 1) * pageSize

	if err := p.db.WithContext(ctx).
		Model(&model.PostTag{}).
		Where("tag_id = ?", tagID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := p.db.WithContext(ctx).
		Model(&model.PostTag{}).
		Where("tag_id = ?", tagID).
		Offset(offset).
		Limit(pageSize).
		Pluck("post_id", &postIDs).Error
	return postIDs, total, err
}

var _ PostTagStore = (*postTagStore)(nil)
