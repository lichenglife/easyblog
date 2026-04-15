package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	postv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/post"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// PostHandler 定义 PostHandler 接口
type PostHandler interface {
	// CreatePost 实现创建帖子接口
	CreatePost(c *gin.Context)
	// DeletePost 实现删除帖子接口
	DeletePost(c *gin.Context)
	// GetPostByID 实现根据 ID 获取帖子接口
	GetPostByID(c *gin.Context)
	// ListPosts 实现获取帖子列表接口
	ListPosts(c *gin.Context)
	// UpdatePost 实现更新帖子接口
	UpdatePost(c *gin.Context)
	// GetPostsByUserID 实现根据用户 ID 获取帖子列表接口
	GetPostsByUserID(c *gin.Context)
	// SearchPosts 搜索文章
	SearchPosts(c *gin.Context)
	// GetMyPosts 获取我的文章
	GetMyPosts(c *gin.Context)
	// TopPost 实现置顶文章接口
	TopPost(c *gin.Context)
	// UnTopPost 实现取消置顶文章接口
	UnTopPost(c *gin.Context)
}

// postHandler 实现 PostHandler 接口
type postHandler struct {
	logger  *log.Logger
	postBiz postv1.PostBiz
}

// NewPostHandler 创建 PostHandler 实例
func NewPostHandler(logger *log.Logger, postBiz postv1.PostBiz) PostHandler {
	return &postHandler{
		logger:  logger,
		postBiz: postBiz,
	}
}

var _ PostHandler = (*postHandler)(nil)

// CreatePost implements PostHandler.
func (p *postHandler) CreatePost(c *gin.Context) {
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		p.logger.Warn("创建文章参数验证失败", zap.Error(err))
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
		p.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	post, err := p.postBiz.CreatePost(c.Request.Context(), uid, &req)
	if err != nil {
		p.logger.Error("创建文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, gin.H{"postId": post.PostID})
}

// DeletePost implements PostHandler.
func (p *postHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		p.logger.Warn("文章 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		p.logger.Warn("文章 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := p.postBiz.DeletePost(c.Request.Context(), id); err != nil {
		p.logger.Error("删除文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// GetPostByID implements PostHandler.
func (p *postHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		p.logger.Warn("文章 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		p.logger.Warn("文章 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	post, err := p.postBiz.GetPostByID(c.Request.Context(), id)
	if err != nil {
		p.logger.Error("获取文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, post)
}

// GetPostsByUserID implements PostHandler.
func (p *postHandler) GetPostsByUserID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		p.logger.Warn("用户 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			p.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			p.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// 将 userID 转换为 uint
	var uid uint
	if _, err := fmt.Sscanf(userID, "%d", &uid); err != nil {
		p.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	posts, err := p.postBiz.GetPostsByUserID(c.Request.Context(), uid, page, pageSize)
	if err != nil {
		p.logger.Error("获取用户文章列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, posts)
}

// ListPosts implements PostHandler.
func (p *postHandler) ListPosts(c *gin.Context) {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			p.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			p.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// 支持按分类 ID 筛选
	categoryIDStr := c.Query("categoryId")
	if categoryIDStr != "" {
		var categoryID uint
		if _, err := fmt.Sscanf(categoryIDStr, "%d", &categoryID); err != nil {
			p.logger.Warn("分类 ID 参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
		posts, err := p.postBiz.ListByCategoryID(c.Request.Context(), categoryID, page, pageSize)
		if err != nil {
			p.logger.Error("按分类获取文章列表失败", zap.Error(err))
			core.WriteResponse(c, err, nil)
			return
		}
		core.WriteResponse(c, errno.OK, posts)
		return
	}

	// 支持按标签 ID 筛选
	tagIDStr := c.Query("tagId")
	if tagIDStr != "" {
		var tagID uint
		if _, err := fmt.Sscanf(tagIDStr, "%d", &tagID); err != nil {
			p.logger.Warn("标签 ID 参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
		posts, err := p.postBiz.ListByTagID(c.Request.Context(), tagID, page, pageSize)
		if err != nil {
			p.logger.Error("按标签获取文章列表失败", zap.Error(err))
			core.WriteResponse(c, err, nil)
			return
		}
		core.WriteResponse(c, errno.OK, posts)
		return
	}

	posts, err := p.postBiz.ListPosts(c.Request.Context(), page, pageSize)
	if err != nil {
		p.logger.Error("获取文章列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, posts)
}

// UpdatePost implements PostHandler.
func (p *postHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		p.logger.Warn("文章 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		p.logger.Warn("文章 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		p.logger.Warn("更新文章参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	// 确保 ID 一致
	req.ID = id

	if err := p.postBiz.UpdatePost(c.Request.Context(), id, &req); err != nil {
		p.logger.Error("更新文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// TopPost 置顶文章
func (p *postHandler) TopPost(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		p.logger.Warn("文章 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		p.logger.Warn("文章 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := p.postBiz.TopPost(c.Request.Context(), id); err != nil {
		p.logger.Error("置顶文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// UnTopPost 取消置顶文章
func (p *postHandler) UnTopPost(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		p.logger.Warn("文章 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		p.logger.Warn("文章 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := p.postBiz.UnTopPost(c.Request.Context(), id); err != nil {
		p.logger.Error("取消置顶文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// SearchPosts 搜索文章
func (p *postHandler) SearchPosts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		p.logger.Warn("搜索关键词不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			p.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			p.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if pageSize > 100 {
		pageSize = 100
	}

	posts, err := p.postBiz.SearchPosts(c.Request.Context(), keyword, page, pageSize)
	if err != nil {
		p.logger.Error("搜索文章失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, posts)
}

// GetMyPosts 获取我的文章列表
func (p *postHandler) GetMyPosts(c *gin.Context) {
	// 从上下文获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			p.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			p.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if pageSize > 100 {
		pageSize = 100
	}

	// 将 userID 转换为 uint
	var uid uint
	if _, err := fmt.Sscanf(userID.(string), "%d", &uid); err != nil {
		p.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	posts, err := p.postBiz.GetMyPosts(c.Request.Context(), uid, page, pageSize)
	if err != nil {
		p.logger.Error("获取我的文章列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, posts)
}
