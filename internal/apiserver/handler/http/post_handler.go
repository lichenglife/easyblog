package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// PostHandler 定义PostHandler接口
type PostHandler interface {
	// createPost 实现创建帖子接口
	CreatePost(c *gin.Context)
	// deletePost 实现删除帖子接口
	DeletePost(c *gin.Context)
	// getPostByID 实现根据ID获取帖子接口
	GetPostByID(c *gin.Context)
	// listPosts 实现获取帖子列表接口
	ListPosts(c *gin.Context)
	// updatePost 实现更新帖子接口
	UpdatePost(c *gin.Context)
	// getPostsByUserID 实现根据用户ID获取帖子列表接口
	GetPostsByUserID(c *gin.Context)
}

// postHandler 实现PostHandler接口
type postHandler struct {
	logger  *log.Logger
	postBiz biz.PostBiz
}

// createPost implements PostHandler.
func (p *postHandler) CreatePost(c *gin.Context) {
	panic("unimplemented")
}

// deletePost implements PostHandler.
func (p *postHandler) DeletePost(c *gin.Context) {
	panic("unimplemented")
}

// getPostByID implements PostHandler.
func (p *postHandler) GetPostByID(c *gin.Context) {
	panic("unimplemented")
}

// getPostsByUserID implements PostHandler.
func (p *postHandler) GetPostsByUserID(c *gin.Context) {
	panic("unimplemented")
}

// listPosts implements PostHandler.
func (p *postHandler) ListPosts(c *gin.Context) {
	panic("unimplemented")
}

// updatePost implements PostHandler.
func (p *postHandler) UpdatePost(c *gin.Context) {
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
