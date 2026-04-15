package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	commentv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/comment"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// CommentHandler 定义 CommentHandler 接口
type CommentHandler interface {
	// CreateComment 创建评论
	CreateComment(c *gin.Context)
	// GetCommentByID 根据 ID 获取评论
	GetCommentByID(c *gin.Context)
	// UpdateComment 更新评论
	UpdateComment(c *gin.Context)
	// DeleteComment 删除评论
	DeleteComment(c *gin.Context)
	// ListComments 获取文章评论列表
	ListComments(c *gin.Context)
	// ListReplies 获取楼中楼回复列表
	ListReplies(c *gin.Context)
	// LikeComment 点赞评论
	LikeComment(c *gin.Context)
	// UnlikeComment 取消点赞评论
	UnlikeComment(c *gin.Context)
}

// commentHandler 实现 CommentHandler 接口
type commentHandler struct {
	logger    *log.Logger
	commentBiz commentv1.CommentBiz
}

// NewCommentHandler 创建 CommentHandler 实例
func NewCommentHandler(logger *log.Logger, commentBiz commentv1.CommentBiz) CommentHandler {
	return &commentHandler{
		logger:    logger,
		commentBiz: commentBiz,
	}
}

var _ CommentHandler = (*commentHandler)(nil)

// CreateComment implements CommentHandler.
func (h *commentHandler) CreateComment(c *gin.Context) {
	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("创建评论参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	// 将 userID 转换为 uint
	var uid uint
	if _, err := fmt.Sscanf(userID.(string), "%d", &uid); err != nil {
		h.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	comment, err := h.commentBiz.CreateComment(c.Request.Context(), uid, &req)
	if err != nil {
		h.logger.Error("创建评论失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, comment)
}

// GetCommentByID implements CommentHandler.
func (h *commentHandler) GetCommentByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Warn("评论 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Warn("评论 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	comment, err := h.commentBiz.GetCommentByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("获取评论失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, comment)
}

// UpdateComment implements CommentHandler.
func (h *commentHandler) UpdateComment(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Warn("评论 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Warn("评论 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var req model.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("更新评论参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	// 确保 ID 一致
	req.ID = id

	// 从上下文获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	// 将 userID 转换为 uint
	var uid uint
	if _, err := fmt.Sscanf(userID.(string), "%d", &uid); err != nil {
		h.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	comment, err := h.commentBiz.UpdateComment(c.Request.Context(), uid, &req)
	if err != nil {
		h.logger.Error("更新评论失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, comment)
}

// DeleteComment implements CommentHandler.
func (h *commentHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Warn("评论 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Warn("评论 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := h.commentBiz.DeleteComment(c.Request.Context(), id); err != nil {
		h.logger.Error("删除评论失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// ListComments implements CommentHandler.
func (h *commentHandler) ListComments(c *gin.Context) {
	postIDStr := c.Param("postID")
	if postIDStr == "" {
		h.logger.Warn("文章 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var postID uint
	if _, err := fmt.Sscanf(postIDStr, "%d", &postID); err != nil {
		h.logger.Warn("文章 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			h.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			h.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if pageSize > 100 {
		pageSize = 100
	}

	comments, err := h.commentBiz.ListComments(c.Request.Context(), postID, page, pageSize)
	if err != nil {
		h.logger.Error("获取评论列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, comments)
}

// ListReplies implements CommentHandler.
func (h *commentHandler) ListReplies(c *gin.Context) {
	parentIDStr := c.Param("parentID")
	if parentIDStr == "" {
		h.logger.Warn("父评论 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var parentID uint
	if _, err := fmt.Sscanf(parentIDStr, "%d", &parentID); err != nil {
		h.logger.Warn("父评论 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			h.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			h.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if pageSize > 100 {
		pageSize = 100
	}

	replies, err := h.commentBiz.ListReplies(c.Request.Context(), parentID, page, pageSize)
	if err != nil {
		h.logger.Error("获取回复列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, replies)
}

// LikeComment implements CommentHandler.
func (h *commentHandler) LikeComment(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Warn("评论 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Warn("评论 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	// 将 userID 转换为 uint
	var uid uint
	if _, err := fmt.Sscanf(userID.(string), "%d", &uid); err != nil {
		h.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := h.commentBiz.LikeComment(c.Request.Context(), uid, id); err != nil {
		h.logger.Error("点赞评论失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// UnlikeComment implements CommentHandler.
func (h *commentHandler) UnlikeComment(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Warn("评论 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		h.logger.Warn("评论 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	// 从上下文获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	// 将 userID 转换为 uint
	var uid uint
	if _, err := fmt.Sscanf(userID.(string), "%d", &uid); err != nil {
		h.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := h.commentBiz.UnlikeComment(c.Request.Context(), uid, id); err != nil {
		h.logger.Error("取消点赞评论失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}
