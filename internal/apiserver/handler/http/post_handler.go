package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// PostHandler 定义PostHandler接口
type PostHandler interface {
	// createPost 实现创建帖子接口
	CreatePost() gin.HandlerFunc
	// deletePost 实现删除帖子接口
	DeletePost() gin.HandlerFunc
	// getPostByID 实现根据ID获取帖子接口
	GetPostByID() gin.HandlerFunc
	// getPostByPostID 实现根据帖子ID获取帖子接口
	GetPostByPostID() gin.HandlerFunc
	// listPosts 实现获取帖子列表接口
	ListPosts() gin.HandlerFunc
	// updatePost 实现更新帖子接口
	UpdatePost() gin.HandlerFunc
	// getPostsByUserID 实现根据用户ID获取帖子列表接口
	GetPostsByUserID() gin.HandlerFunc
}

// postHandler 实现PostHandler接口
type postHandler struct {
	logger  *log.Logger
	postBiz biz.PostBiz
}

// createPost implements PostHandler.
func (p *postHandler) CreatePost() gin.HandlerFunc {
	panic("unimplemented")
}

// deletePost implements PostHandler.
func (p *postHandler) DeletePost() gin.HandlerFunc {
	panic("unimplemented")
}

// getPostByID implements PostHandler.
func (p *postHandler) GetPostByID() gin.HandlerFunc {
	panic("unimplemented")
}

// getPostByPostID implements PostHandler.
func (p *postHandler) GetPostByPostID() gin.HandlerFunc {
	panic("unimplemented")
}

// getPostsByUserID implements PostHandler.
func (p *postHandler) GetPostsByUserID() gin.HandlerFunc {
	panic("unimplemented")
}

// listPosts implements PostHandler.
func (p *postHandler) ListPosts() gin.HandlerFunc {
	panic("unimplemented")
}

// updatePost implements PostHandler.
func (p *postHandler) UpdatePost() gin.HandlerFunc {
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
