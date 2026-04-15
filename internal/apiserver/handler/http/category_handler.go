package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	categoryv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/category"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// CategoryHandler 分类管理接口
type CategoryHandler interface {
	// CreateCategory 创建分类
	CreateCategory(c *gin.Context)
	// UpdateCategory 更新分类
	UpdateCategory(c *gin.Context)
	// DeleteCategory 删除分类
	DeleteCategory(c *gin.Context)
	// GetCategory 获取分类详情
	GetCategory(c *gin.Context)
	// GetCategoryTree 获取分类树
	GetCategoryTree(c *gin.Context)
	// ListCategories 获取分类列表
	ListCategories(c *gin.Context)
}

// categoryHandler 实现 CategoryHandler 接口
type categoryHandler struct {
	logger      *log.Logger
	categoryBiz categoryv1.CategoryBiz
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(logger *log.Logger, categoryBiz categoryv1.CategoryBiz) CategoryHandler {
	return &categoryHandler{
		logger:      logger,
		categoryBiz: categoryBiz,
	}
}

var _ CategoryHandler = (*categoryHandler)(nil)

// CreateCategory 创建分类
func (ch *categoryHandler) CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ch.logger.Warn("创建分类参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	category, err := ch.categoryBiz.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		ch.logger.Error("创建分类失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, gin.H{"categoryId": category.ID})
}

// UpdateCategory 更新分类
func (ch *categoryHandler) UpdateCategory(c *gin.Context) {
	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ch.logger.Warn("更新分类参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	category, err := ch.categoryBiz.UpdateCategory(c.Request.Context(), &req)
	if err != nil {
		ch.logger.Error("更新分类失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, gin.H{"categoryId": category.ID})
}

// DeleteCategory 删除分类
func (ch *categoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		ch.logger.Warn("分类 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		ch.logger.Warn("分类 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := ch.categoryBiz.DeleteCategory(c.Request.Context(), id); err != nil {
		ch.logger.Error("删除分类失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// GetCategory 获取分类详情
func (ch *categoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		ch.logger.Warn("分类 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		ch.logger.Warn("分类 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	category, err := ch.categoryBiz.GetCategory(c.Request.Context(), id)
	if err != nil {
		ch.logger.Error("获取分类失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, category)
}

// GetCategoryTree 获取分类树
func (ch *categoryHandler) GetCategoryTree(c *gin.Context) {
	categories, err := ch.categoryBiz.GetCategoryTree(c.Request.Context())
	if err != nil {
		ch.logger.Error("获取分类树失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, gin.H{"categories": categories})
}

// ListCategories 获取分类列表
func (ch *categoryHandler) ListCategories(c *gin.Context) {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			ch.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			ch.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if pageSize > 100 {
		pageSize = 100
	}

	categories, total, err := ch.categoryBiz.ListCategories(c.Request.Context(), page, pageSize)
	if err != nil {
		ch.logger.Error("获取分类列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, gin.H{
		"totalCount": total,
		"hasMore":    len(categories) == pageSize,
		"categories": categories,
	})
}
