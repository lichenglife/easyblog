package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

// TagStore 标签存储接口
type TagStore interface {
	Create(ctx context.Context, tag *model.Tag) error
	Update(ctx context.Context, tag *model.Tag) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*model.Tag, error)
	GetByName(ctx context.Context, name string) (*model.Tag, error)
	List(ctx context.Context, offset, limit int) ([]*model.Tag, int64, error)
	ListAll(ctx context.Context) ([]*model.Tag, error)
	UpdatePostCount(ctx context.Context, id uint, delta int) error
}

// tagStore 标签存储实现
type tagStore struct {
	db *gorm.DB
}

// 确保 tagStore 实现 TagStore 接口
var _ TagStore = (*tagStore)(nil)

// NewTags 创建标签存储实例
func NewTags(db *gorm.DB) TagStore {
	return &tagStore{db: db}
}

// Create 创建标签
func (ts *tagStore) Create(ctx context.Context, tag *model.Tag) error {
	return ts.db.WithContext(ctx).Create(tag).Error
}

// Update 更新标签
func (ts *tagStore) Update(ctx context.Context, tag *model.Tag) error {
	return ts.db.WithContext(ctx).Save(tag).Error
}

// Delete 删除标签
func (ts *tagStore) Delete(ctx context.Context, id uint) error {
	return ts.db.WithContext(ctx).Delete(&model.Tag{}, id).Error
}

// GetByID 根据 ID 获取标签
func (ts *tagStore) GetByID(ctx context.Context, id uint) (*model.Tag, error) {
	var tag model.Tag
	err := ts.db.WithContext(ctx).First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetBySlug 根据别名获取标签
func (ts *tagStore) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	var tag model.Tag
	err := ts.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetByName 根据名称获取标签
func (ts *tagStore) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	var tag model.Tag
	err := ts.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// List 获取标签列表（分页）
func (ts *tagStore) List(ctx context.Context, offset, limit int) ([]*model.Tag, int64, error) {
	var tags []*model.Tag
	var total int64

	if err := ts.db.WithContext(ctx).Model(&model.Tag{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := ts.db.WithContext(ctx).
		Order("post_count DESC, id DESC").
		Offset(offset).
		Limit(limit).
		Find(&tags).Error

	return tags, total, err
}

// ListAll 获取所有标签
func (ts *tagStore) ListAll(ctx context.Context) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := ts.db.WithContext(ctx).
		Order("post_count DESC, id DESC").
		Find(&tags).Error
	return tags, err
}

// UpdatePostCount 更新标签文章数（增加或减少）
func (ts *tagStore) UpdatePostCount(ctx context.Context, id uint, delta int) error {
	return ts.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var tag model.Tag
		if err := tx.First(&tag, id).Error; err != nil {
			return err
		}

		newCount := int(tag.PostCount) + delta
		if newCount < 0 {
			return nil // 文章数不能为负，但不返回错误
		}

		return tx.Model(&tag).Update("post_count", uint(newCount)).Error
	})
}
