package handler

import (
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/auth"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// 定义 Handler 接口
type Handler interface {

	// User 用户相关接口
	Users() UserHandler

	// Post 帖子相关接口
	Posts() PostHandler

	// Comment 评论相关接口
	Comments() CommentHandler

	// Like 点赞相关接口
	Likes() LikeHandler

	// Category 分类相关接口
	Categories() CategoryHandler

	// Tag 标签相关接口
	Tags() TagHandler

	// Image 图片相关接口
	Images() ImageHandler
}

// handler 定义了 Handler 接口的实现
type handler struct {
	logger          *log.Logger
	store           store.IStore
	UserHandler     UserHandler
	PostHandler     PostHandler
	CommentHandler  CommentHandler
	LikeHandler     LikeHandler
	CategoryHandler CategoryHandler
	TagHandler      TagHandler
	ImageHandler    ImageHandler
}

// NewHandler 创建 Handler 实例
func NewHandler(logger *log.Logger, store store.IStore, jwt *auth.JWT) Handler {
	h := &handler{
		logger: logger,
		store:  store,
	}
	b := biz.NewBiz(store, jwt, logger.Logger)

	h.UserHandler = NewUserHandler(logger, b)
	h.PostHandler = NewPostHandler(logger, b.PostV1())
	h.CommentHandler = NewCommentHandler(logger, b.CommentV1())
	h.LikeHandler = NewLikeHandler(logger, b.LikeV1())
	h.CategoryHandler = NewCategoryHandler(logger, b.CategoryV1())
	h.TagHandler = NewTagHandler(logger, b.TagV1())
	h.ImageHandler = NewImageHandler(logger, "./uploads/images", 5<<20) // 5MB 限制
	return h
}

// User 用户相关接口
func (h *handler) Users() UserHandler {
	return h.UserHandler
}

// Post 帖子相关接口
func (h *handler) Posts() PostHandler {
	return h.PostHandler
}

// Comment 评论相关接口
func (h *handler) Comments() CommentHandler {
	return h.CommentHandler
}

// Like 点赞相关接口
func (h *handler) Likes() LikeHandler {
	return h.LikeHandler
}

// Category 分类相关接口
func (h *handler) Categories() CategoryHandler {
	return h.CategoryHandler
}

// Tag 标签相关接口
func (h *handler) Tags() TagHandler {
	return h.TagHandler
}

// Image 图片相关接口
func (h *handler) Images() ImageHandler {
	return h.ImageHandler
}
