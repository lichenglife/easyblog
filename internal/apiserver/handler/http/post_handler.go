package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// PostHandler 定义PostHandler接口
type PostHandler interface {
	// createPost 实现创建帖子接口
	CreatePost(c *gin.Context)
	// deletePost 实现删除帖子接口
	DeletePost(c *gin.Context)
	// getPostByID 实现根据ID获取帖子接口
	GetPostByID(c *gin.Context)
	// getPostByPostID 实现根据帖子ID获取帖子接口
	GetPostByPostID(c *gin.Context)
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

// @Summary      创建博客文章
// @Description  创建一篇新的博客文章
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        request body model.CreatePostRequest true "创建文章请求参数"
// @Success      200  {object}  model.Post
// @Failure      400  {object}  core.ErrResponse
// @Failure      500  {object}  core.ErrResponse
// @Security     Bearer
// @Router       /v1/posts [post]
// createPost implements PostHandler.
func (p *postHandler) CreatePost(c *gin.Context) {

	userID := c.GetString("userID")
	if userID == "" {
		core.WriteResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	var postRequest model.CreatePostRequest
	if err := c.ShouldBindJSON(&postRequest); err != nil {
		p.logger.Logger.Error("请求参数错误", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	postRequest.UserID = userID
	post, err := p.postBiz.CreatePost(c, &postRequest)
	if err != nil {
		p.logger.Logger.Error("创建博客失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, post)

}

// @Summary      删除博客文章
// @Description  删除指定 ID 的博客文章
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        postID path string true "文章ID"
// @Success      200  {object}  string "博客删除成功"
// @Failure      400  {object}  core.ErrResponse
// @Failure      500  {object}  core.ErrResponse
// @Security     Bearer
// @Router       /v1/posts/{postID} [delete]
// deletePost implements PostHandler.
func (p *postHandler) DeletePost(c *gin.Context) {
	postID := c.Param("postID")
	err := p.postBiz.DeletePostByPostID(c, postID)
	if err != nil {
		p.logger.Logger.Error("删除博客失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)
}

// getPostByID implements PostHandler.
func (p *postHandler) GetPostByID(c *gin.Context) {
	postID := c.Param("postID")
	post, err := p.postBiz.GetPostByPostID(c, postID)
	if err != nil {
		p.logger.Logger.Error("获取博客失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, post)
}

// @Summary      获取博客文章
// @Description  获取指定 ID 的博客文章
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        postID path string true "文章ID"
// @Success      200  {object}  model.Post
// @Failure      400  {object}  core.ErrResponse
// @Failure      500  {object}  core.ErrResponse
// @Security     Bearer
// @Router       /v1/posts/{postID} [get]
// getPostByPostID implements PostHandler.
func (p *postHandler) GetPostByPostID(c *gin.Context) {

	postID := c.Param("postID")
	post, err := p.postBiz.GetPostByPostID(c, postID)
	if err != nil {
		p.logger.Logger.Error("获取博客失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, post)

}

// @Summary      获取用户的博客文章列表
// @Description  分页获取指定用户的博客文章列表
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        userID path string true "用户名"
// @Param        offset query int false "起始位置" default(0)
// @Param        limit query int false "列表数量" default(10)
// @Success      200  {object}  model.ListPostResponse
// @Failure      400  {object}  core.ErrResponse
// @Failure      500  {object}  core.ErrResponse
// @Security     Bearer
// @Router       /v1/posts/user/{userID} [get]
// getPostsByUserID implements PostHandler.
func (p *postHandler) GetPostsByUserID(c *gin.Context) {
	userID := c.Param("userID")

	pagesize := core.GetLimitParam(c)
	page := core.GetPageParam(c)
	_, posts, err := p.postBiz.GetPostsByUserID(c, userID, page, pagesize)
	if err != nil {
		p.logger.Logger.Error("failed get postByUserID", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, posts)
}

// @Summary      获取博客文章列表
// @Description  分页获取博客文章列表
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        page query int false "起始位置" default(0)
// @Param        limit query int false "列表数量" default(10)
// @Success      200  {object}  model.ListPostResponse
// @Failure      400  {object}  core.ErrResponse
// @Failure      500  {object}  core.ErrResponse
// @Security     Bearer
// @Router       /v1/posts [get]
// listPosts implements PostHandler.
func (p *postHandler) ListPosts(c *gin.Context) {
	page := core.GetPageParam(c)
	pagesize := core.GetLimitParam(c)
	posts, err := p.postBiz.ListPosts(c, page, pagesize)
	if err != nil {
		p.logger.Logger.Error("获取博客列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, posts)
}

// @Summary      更新博客文章
// @Description  更新指定 ID 的博客文章
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        postID path string true "文章ID"
// @Param        request body model.UpdatePostRequest true "更新文章请求参数"
// @Success      200  {object}  model.UpdatePostRequest
// @Failure      400  {object}  core.ErrResponse
// @Failure      500  {object}  core.ErrResponse
// @Security     Bearer
// @Router       /v1/posts/{postID} [put]
// updatePost implements PostHandler.
func (p *postHandler) UpdatePost(c *gin.Context) {
	postID := c.Param("postID")
	var postRequest model.UpdatePostRequest
	if err := c.ShouldBindJSON(&postRequest); err != nil {
		p.logger.Logger.Error("请求参数错误", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	postRequest.PostID = postID
	err := p.postBiz.UpdatePost(c, &postRequest)
	if err != nil {
		p.logger.Logger.Error("更新博客失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, postRequest)
}

// NewPostHandler 创建PostHandler实例
func NewPostHandler(logger *log.Logger, postBiz biz.PostBiz) PostHandler {
	return &postHandler{
		logger:  logger,
		postBiz: postBiz,
	}
}

var _ PostHandler = (*postHandler)(nil)
