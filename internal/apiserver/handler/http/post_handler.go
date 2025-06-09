package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
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
	postBiz biz.IBiz
}

// createPost implements PostHandler.
func (p *postHandler) CreatePost(c *gin.Context) {
	// // 1、解析参数
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Log.Error(err.Error())
		core.WriteResponse(c, err, nil)
		return
	}
	// 2、当前用户
	userID := c.GetString("userID")
	req.UserID = userID
	//3、创建博客
	post, err := p.postBiz.PostV1().CreatePost(c, &req)
	if err != nil {
		log.Log.Error(err.Error())
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, post)

}

// deletePost implements PostHandler.
func (p *postHandler) DeletePost(c *gin.Context) {
	// 1、解析参数
	postID := c.Param("postID")
	// 2、删除博客
	if err := p.postBiz.PostV1().DeletePostByPostID(c, postID); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, errno.OK)

}

// getPostByID implements PostHandler.
func (p *postHandler) GetPostByID(c *gin.Context) {
	// 1、解析参数
	postID := c.Param("postID")
	// 2、执行查询
	post, err := p.postBiz.PostV1().GetPostByPostID(c, postID)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, post)
}

// getPostsByUserID implements PostHandler.
func (p *postHandler) GetPostsByUserID(c *gin.Context) {
	// 1、解析参数
	userID := c.Param("userID")

	pageSize := core.GetLimitParam(c)
	page := core.GetPageParam(c)

	// 2、执行查询
	posts, err := p.postBiz.PostV1().GetPostsByUserID(c, userID, page, pageSize)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, posts)

}

// listPosts implements PostHandler.
func (p *postHandler) ListPosts(c *gin.Context) {
	// 1、解析参数

	pageSize := core.GetLimitParam(c)
	page := core.GetPageParam(c)

	// 2、执行查询
	posts, err := p.postBiz.PostV1().ListPosts(c, page, pageSize)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, posts)
}

// updatePost implements PostHandler.
func (p *postHandler) UpdatePost(c *gin.Context) {
	// 1、解析参数
	postID := c.Param("postID")
	var req *model.UpdatePostRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	// 2、执行更新
	req.PostID = postID
	if err := p.postBiz.PostV1().UpdatePost(c, req); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, "更新成功")
}

// NewPostHandler 创建PostHandler实例
func NewPostHandler(logger *log.Logger, postBiz biz.IBiz) PostHandler {
	return &postHandler{
		postBiz: postBiz,
	}
}

var _ PostHandler = (*postHandler)(nil)
