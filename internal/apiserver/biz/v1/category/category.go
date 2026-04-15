package category

import (
	"context"
	"fmt"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"go.uber.org/zap"
)

// CategoryBiz 分类业务接口
type CategoryBiz interface {
	CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.Category, error)
	UpdateCategory(ctx context.Context, req *model.UpdateCategoryRequest) (*model.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
	GetCategory(ctx context.Context, id uint) (*model.Category, error)
	GetCategoryTree(ctx context.Context) ([]*model.Category, error)
	ListCategories(ctx context.Context, page, pageSize int) ([]*model.Category, int64, error)
}

// categoryBiz 分类业务实现
type categoryBiz struct {
	logger *zap.Logger
	store  store.IStore
}

// 确保 categoryBiz 实现 CategoryBiz 接口
var _ CategoryBiz = (*categoryBiz)(nil)

// NewCategoryBiz 创建分类业务实例
func NewCategoryBiz(logger *zap.Logger, storeInstance store.IStore) CategoryBiz {
	return &categoryBiz{
		logger: logger,
		store:  storeInstance,
	}
}

// CreateCategory 创建分类
func (cb *categoryBiz) CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.Category, error) {
	// 检查名称是否已存在
	existing, err := cb.store.Category().GetBySlug(ctx, req.Slug)
	if err == nil && existing != nil {
		return nil, errno.ErrCategoryAlreadyExist
	}

	// 如果指定了父分类，验证父分类存在
	var level int = 1
	if req.ParentID > 0 {
		parent, err := cb.store.Category().GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("父分类不存在")
		}
		level = parent.Level + 1
		if level > 3 {
			return nil, fmt.Errorf("分类层级最多支持 3 级")
		}
	}

	category := &model.Category{
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: req.ParentID,
		Level:    level,
		Sort:     req.Sort,
		Status:   1,
	}

	if err := cb.store.Category().Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory 更新分类
func (cb *categoryBiz) UpdateCategory(ctx context.Context, req *model.UpdateCategoryRequest) (*model.Category, error) {
	// 获取原分类
	category, err := cb.store.Category().GetByID(ctx, req.ID)
	if err != nil {
		return nil, errno.ErrCategoryNotFound
	}

	// 检查新别名是否已被其他分类使用
	if req.Slug != category.Slug {
		existing, err := cb.store.Category().GetBySlug(ctx, req.Slug)
		if err == nil && existing != nil && existing.ID != req.ID {
			return nil, errno.ErrCategoryAlreadyExist
		}
	}

	// 如果修改了父分类，验证层级
	if req.ParentID != category.ParentID && req.ParentID > 0 {
		parent, err := cb.store.Category().GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("父分类不存在")
		}
		// 不能将自己设为父分类
		if parent.ID == req.ID {
			return nil, fmt.Errorf("不能将自己设为父分类")
		}
		category.Level = parent.Level + 1
		if category.Level > 3 {
			return nil, fmt.Errorf("分类层级最多支持 3 级")
		}
	}

	// 更新字段
	category.Name = req.Name
	category.Slug = req.Slug
	category.ParentID = req.ParentID
	category.Sort = req.Sort
	if req.Status != 0 {
		category.Status = req.Status
	}

	if err := cb.store.Category().Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory 删除分类
func (cb *categoryBiz) DeleteCategory(ctx context.Context, id uint) error {
	category, err := cb.store.Category().GetByID(ctx, id)
	if err != nil {
		return errno.ErrCategoryNotFound
	}

	// 检查是否有子分类
	children, err := cb.store.Category().ListAll(ctx)
	if err != nil {
		return err
	}
	for _, child := range children {
		if child.ParentID == id {
			return fmt.Errorf("存在子分类，无法删除")
		}
	}

	// 检查是否有文章关联
	if category.PostCount > 0 {
		return fmt.Errorf("分类下有文章，无法删除")
	}

	return cb.store.Category().Delete(ctx, id)
}

// GetCategory 获取分类详情
func (cb *categoryBiz) GetCategory(ctx context.Context, id uint) (*model.Category, error) {
	category, err := cb.store.Category().GetByID(ctx, id)
	if err != nil {
		return nil, errno.ErrCategoryNotFound
	}
	return category, nil
}

// GetCategoryTree 获取分类树
func (cb *categoryBiz) GetCategoryTree(ctx context.Context) ([]*model.Category, error) {
	categories, err := cb.store.Category().GetTree(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// ListCategories 获取分类列表（分页）
func (cb *categoryBiz) ListCategories(ctx context.Context, page, pageSize int) ([]*model.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return cb.store.Category().List(ctx, offset, pageSize)
}
