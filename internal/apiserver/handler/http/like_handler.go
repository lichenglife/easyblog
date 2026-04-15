package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	likev1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/like"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// LikeHandler 定义 LikeHandler 接口
type LikeHandler interface {
	// LikePost 点赞文章
	LikePost(c *gin.Context)
	// UnlikePost 取消点赞文章
	UnlikePost(c *gin.Context)
	// GetLikeStatus 获取点赞状态
	GetLikeStatus(c *gin.Context)
}

// likeHandler 实现 LikeHandler 接口
type likeHandler struct {
	logger  *log.Logger
	likeBiz likev1.LikeBiz
}

// NewLikeHandler 创建 LikeHandler 实例
func NewLikeHandler(logger *log.Logger, likeBiz likev1.LikeBiz) LikeHandler {
	return &likeHandler{
		logger:  logger,
		likeBiz: likeBiz,
	}
}

var _ LikeHandler = (*likeHandler)(nil)

// LikePost implements LikeHandler.
func (h *likeHandler) LikePost(c *gin.Context) {
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

	if err := h.likeBiz.LikePost(c.Request.Context(), uid, postID); err != nil {
		h.logger.Error("点赞文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// UnlikePost implements LikeHandler.
func (h *likeHandler) UnlikePost(c *gin.Context) {
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

	if err := h.likeBiz.UnlikePost(c.Request.Context(), uid, postID); err != nil {
		h.logger.Error("取消点赞文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// GetLikeStatus implements LikeHandler.
func (h *likeHandler) GetLikeStatus(c *gin.Context) {
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

	// 从上下文获取用户 ID（可选，用于返回当前用户的点赞状态）
	userID, _ := c.Get("userID")

	var uid uint
	if userID != nil {
		if _, err := fmt.Sscanf(userID.(string), "%d", &uid); err != nil {
			h.logger.Warn("用户 ID 格式错误", zap.Error(err))
		}
	}

	isLiked := false
	if uid > 0 {
		isLiked, _ = h.likeBiz.IsLikedPost(c.Request.Context(), uid, postID)
	}

	// 获取点赞总数
	count, _ := h.likeBiz.GetPostLikeCount(c.Request.Context(), postID)

	core.WriteResponse(c, errno.OK, gin.H{
		"isLiked": isLiked,
		"count":   count,
	})
}
