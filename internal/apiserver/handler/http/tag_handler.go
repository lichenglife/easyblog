package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	tagv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/tag"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// TagHandler 标签管理接口
type TagHandler interface {
	// CreateTag 创建标签
	CreateTag(c *gin.Context)
	// UpdateTag 更新标签
	UpdateTag(c *gin.Context)
	// DeleteTag 删除标签
	DeleteTag(c *gin.Context)
	// GetTag 获取标签详情
	GetTag(c *gin.Context)
	// ListTags 获取标签列表
	ListTags(c *gin.Context)
}

// tagHandler 实现 TagHandler 接口
type tagHandler struct {
	logger *log.Logger
	tagBiz tagv1.TagBiz
}

// NewTagHandler 创建标签处理器
func NewTagHandler(logger *log.Logger, tagBiz tagv1.TagBiz) TagHandler {
	return &tagHandler{
		logger: logger,
		tagBiz: tagBiz,
	}
}

var _ TagHandler = (*tagHandler)(nil)

// CreateTag 创建标签
func (th *tagHandler) CreateTag(c *gin.Context) {
	var req model.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		th.logger.Warn("创建标签参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	tag, err := th.tagBiz.CreateTag(c.Request.Context(), &req)
	if err != nil {
		th.logger.Error("创建标签失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, gin.H{"tagId": tag.ID})
}

// UpdateTag 更新标签
func (th *tagHandler) UpdateTag(c *gin.Context) {
	var req model.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		th.logger.Warn("更新标签参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	tag, err := th.tagBiz.UpdateTag(c.Request.Context(), &req)
	if err != nil {
		th.logger.Error("更新标签失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, gin.H{"tagId": tag.ID})
}

// DeleteTag 删除标签
func (th *tagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		th.logger.Warn("标签 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		th.logger.Warn("标签 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := th.tagBiz.DeleteTag(c.Request.Context(), id); err != nil {
		th.logger.Error("删除标签失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// GetTag 获取标签详情
func (th *tagHandler) GetTag(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		th.logger.Warn("标签 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		th.logger.Warn("标签 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	tag, err := th.tagBiz.GetTag(c.Request.Context(), id)
	if err != nil {
		th.logger.Error("获取标签失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, tag)
}

// ListTags 获取标签列表
func (th *tagHandler) ListTags(c *gin.Context) {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			th.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			th.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if pageSize > 100 {
		pageSize = 100
	}

	tags, total, err := th.tagBiz.ListTags(c.Request.Context(), page, pageSize)
	if err != nil {
		th.logger.Error("获取标签列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, gin.H{
		"totalCount": total,
		"hasMore":    len(tags) == pageSize,
		"tags":       tags,
	})
}
