package tag

import (
	"context"
	"fmt"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"go.uber.org/zap"
)

// TagBiz 标签业务接口
type TagBiz interface {
	CreateTag(ctx context.Context, req *model.CreateTagRequest) (*model.Tag, error)
	UpdateTag(ctx context.Context, req *model.UpdateTagRequest) (*model.Tag, error)
	DeleteTag(ctx context.Context, id uint) error
	GetTag(ctx context.Context, id uint) (*model.Tag, error)
	ListTags(ctx context.Context, page, pageSize int) ([]*model.Tag, int64, error)
}

// tagBiz 标签业务实现
type tagBiz struct {
	logger *zap.Logger
	store  store.IStore
}

// 确保 tagBiz 实现 TagBiz 接口
var _ TagBiz = (*tagBiz)(nil)

// NewTagBiz 创建标签业务实例
func NewTagBiz(logger *zap.Logger, storeInstance store.IStore) TagBiz {
	return &tagBiz{
		logger: logger,
		store:  storeInstance,
	}
}

// CreateTag 创建标签
func (tb *tagBiz) CreateTag(ctx context.Context, req *model.CreateTagRequest) (*model.Tag, error) {
	// 检查名称是否已存在
	existing, err := tb.store.Tag().GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, errno.ErrTagAlreadyExist
	}

	// 检查别名是否已存在
	existing, err = tb.store.Tag().GetBySlug(ctx, req.Slug)
	if err == nil && existing != nil {
		return nil, errno.ErrTagAlreadyExist
	}

	tag := &model.Tag{
		Name:   req.Name,
		Slug:   req.Slug,
		Status: 1,
	}

	if err := tb.store.Tag().Create(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// UpdateTag 更新标签
func (tb *tagBiz) UpdateTag(ctx context.Context, req *model.UpdateTagRequest) (*model.Tag, error) {
	// 获取原标签
	tag, err := tb.store.Tag().GetByID(ctx, req.ID)
	if err != nil {
		return nil, errno.ErrTagNotFound
	}

	// 检查新名称是否已被其他标签使用
	if req.Name != tag.Name {
		existing, err := tb.store.Tag().GetByName(ctx, req.Name)
		if err == nil && existing != nil && existing.ID != req.ID {
			return nil, errno.ErrTagAlreadyExist
		}
	}

	// 检查新别名是否已被其他标签使用
	if req.Slug != tag.Slug {
		existing, err := tb.store.Tag().GetBySlug(ctx, req.Slug)
		if err == nil && existing != nil && existing.ID != req.ID {
			return nil, errno.ErrTagAlreadyExist
		}
	}

	// 更新字段
	tag.Name = req.Name
	tag.Slug = req.Slug
	if req.Status != 0 {
		tag.Status = req.Status
	}

	if err := tb.store.Tag().Update(ctx, tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag 删除标签
func (tb *tagBiz) DeleteTag(ctx context.Context, id uint) error {
	tag, err := tb.store.Tag().GetByID(ctx, id)
	if err != nil {
		return errno.ErrTagNotFound
	}

	// 检查是否有文章关联
	if tag.PostCount > 0 {
		return fmt.Errorf("标签下有文章，无法删除")
	}

	return tb.store.Tag().Delete(ctx, id)
}

// GetTag 获取标签详情
func (tb *tagBiz) GetTag(ctx context.Context, id uint) (*model.Tag, error) {
	tag, err := tb.store.Tag().GetByID(ctx, id)
	if err != nil {
		return nil, errno.ErrTagNotFound
	}
	return tag, nil
}

// ListTags 获取标签列表（分页）
func (tb *tagBiz) ListTags(ctx context.Context, page, pageSize int) ([]*model.Tag, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return tb.store.Tag().List(ctx, offset, pageSize)
}
