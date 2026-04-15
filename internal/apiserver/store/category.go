package store

import (
	"context"
	"fmt"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

// CategoryStore 分类存储接口
type CategoryStore interface {
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.Category, error)
	GetBySlug(ctx context.Context, slug string) (*model.Category, error)
	List(ctx context.Context, offset, limit int) ([]*model.Category, int64, error)
	ListAll(ctx context.Context) ([]*model.Category, error)
	GetTree(ctx context.Context) ([]*model.Category, error)
	UpdatePostCount(ctx context.Context, id uint, delta int) error
}

// categoryStore 分类存储实现
type categoryStore struct {
	db *gorm.DB
}

// 确保 categoryStore 实现 CategoryStore 接口
var _ CategoryStore = (*categoryStore)(nil)

// NewCategories 创建分类存储实例
func NewCategories(db *gorm.DB) CategoryStore {
	return &categoryStore{db: db}
}

// Create 创建分类
func (cs *categoryStore) Create(ctx context.Context, category *model.Category) error {
	return cs.db.WithContext(ctx).Create(category).Error
}

// Update 更新分类
func (cs *categoryStore) Update(ctx context.Context, category *model.Category) error {
	return cs.db.WithContext(ctx).Save(category).Error
}

// Delete 删除分类
func (cs *categoryStore) Delete(ctx context.Context, id uint) error {
	return cs.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

// GetByID 根据 ID 获取分类
func (cs *categoryStore) GetByID(ctx context.Context, id uint) (*model.Category, error) {
	var category model.Category
	err := cs.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetBySlug 根据别名获取分类
func (cs *categoryStore) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	var category model.Category
	err := cs.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// List 获取分类列表（分页）
func (cs *categoryStore) List(ctx context.Context, offset, limit int) ([]*model.Category, int64, error) {
	var categories []*model.Category
	var total int64

	if err := cs.db.WithContext(ctx).Model(&model.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := cs.db.WithContext(ctx).
		Order("sort ASC, id DESC").
		Offset(offset).
		Limit(limit).
		Find(&categories).Error

	return categories, total, err
}

// ListAll 获取所有分类
func (cs *categoryStore) ListAll(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	err := cs.db.WithContext(ctx).
		Order("sort ASC, id DESC").
		Find(&categories).Error
	return categories, err
}

// GetTree 获取分类树（按层级排序）
func (cs *categoryStore) GetTree(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	err := cs.db.WithContext(ctx).
		Order("level ASC, sort ASC, id ASC").
		Find(&categories).Error
	return categories, err
}

// UpdatePostCount 更新分类文章数（增加或减少）
func (cs *categoryStore) UpdatePostCount(ctx context.Context, id uint, delta int) error {
	return cs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var category model.Category
		if err := tx.First(&category, id).Error; err != nil {
			return err
		}

		newCount := int(category.PostCount) + delta
		if newCount < 0 {
			return fmt.Errorf("文章数不能为负数")
		}

		return tx.Model(&category).Update("post_count", uint(newCount)).Error
	})
}
