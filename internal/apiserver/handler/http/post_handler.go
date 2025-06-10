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

// CreatePost 创建帖子
// @Summary 创建帖子
// @Description 创建一个新的帖子
// @Tags 帖子管理
// @Accept json
// @Produce json
// @Security     Bearer
// @Param body body model.CreatePostRequest true "创建帖子请求参数"
// @Success 200 {object} core.Response{data=model.Post} "成功响应"
// @Failure 400 {object} core.Response "请求参数错误"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /posts [post]
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

// DeletePost 删除帖子
// @Summary 删除帖子
// @Description 根据帖子 ID 删除帖子
// @Tags 帖子管理
// @Security     Bearer
// @Param postID path string true "帖子 ID"
// @Success 200 {object} core.Response "成功响应"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /posts/{postID} [delete]
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

// GetPostByID 获取帖子信息
// @Summary 获取帖子信息
// @Description 根据帖子 ID 获取帖子信息
// @Tags 帖子管理
// @Security     Bearer
// @Param postID path string true "帖子 ID"
// @Success 200 {object} core.Response{data=model.Post} "成功响应"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /posts/{postID} [get]
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

// GetPostsByUserID 获取用户的帖子列表
// @Summary 获取用户的帖子列表
// @Description 根据用户 ID 获取帖子列表
// @Tags 帖子管理
// @Security     Bearer
// @Param userID path string true "用户 ID"
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} core.Response{data=[]model.Post} "成功响应"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /posts/user/{userID} [get]
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

// ListPosts 获取帖子列表
// @Summary 获取帖子列表
// @Description 分页获取帖子列表
// @Tags 帖子管理
// @Security     Bearer
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} core.Response{data=[]model.ListPostResponse} "成功响应"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /posts [get]
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

// UpdatePost 更新帖子
// @Summary 更新帖子
// @Description 根据帖子 ID 更新帖子
// @Tags 帖子管理
// @Security     Bearer
// @Param postID path string true "帖子 ID"
// @Param body body model.UpdatePostRequest true "更新帖子请求参数"
// @Success 200 {object} core.Response "成功响应"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /posts/{postID} [put]
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
