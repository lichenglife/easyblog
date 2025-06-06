package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// PostHandler 定义PostHandler接口
type PostHandler interface {
	// createPost 实现创建帖子接口
	createPost() gin.HandlerFunc
	// deletePost 实现删除帖子接口
	deletePost() gin.HandlerFunc
	// getPostByID 实现根据ID获取帖子接口
	getPostByID() gin.HandlerFunc
	// getPostByPostID 实现根据帖子ID获取帖子接口
	getPostByPostID() gin.HandlerFunc
	// listPosts 实现获取帖子列表接口
	listPosts() gin.HandlerFunc
	// updatePost 实现更新帖子接口
	updatePost() gin.HandlerFunc
	// getPostsByUserID 实现根据用户ID获取帖子列表接口
	getPostsByUserID() gin.HandlerFunc
}

// postHandler 实现PostHandler接口
type postHandler struct {
	logger  *log.Logger
	postBiz biz.PostBiz
}

// createPost implements PostHandler.
func (p *postHandler) createPost() gin.HandlerFunc {
	panic("unimplemented")
}

// deletePost implements PostHandler.
func (p *postHandler) deletePost() gin.HandlerFunc {
	panic("unimplemented")
}

// getPostByID implements PostHandler.
func (p *postHandler) getPostByID() gin.HandlerFunc {
	panic("unimplemented")
}

// getPostByPostID implements PostHandler.
func (p *postHandler) getPostByPostID() gin.HandlerFunc {
	panic("unimplemented")
}

// getPostsByUserID implements PostHandler.
func (p *postHandler) getPostsByUserID() gin.HandlerFunc {
	panic("unimplemented")
}

// listPosts implements PostHandler.
func (p *postHandler) listPosts() gin.HandlerFunc {
	panic("unimplemented")
}

// updatePost implements PostHandler.
func (p *postHandler) updatePost() gin.HandlerFunc {
	panic("unimplemented")
}

// NewPostHandler 创建PostHandler实例
func NewPostHandler(logger *log.Logger, postBiz biz.PostBiz) PostHandler {
	return &postHandler{
		logger:  logger,
		postBiz: postBiz,
	}
}

var _ PostHandler = (*postHandler)(nil)
