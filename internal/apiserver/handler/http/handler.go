package handler

import (
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/lichenglife/easyblog/internal/pkg/middleware"
)

// 定义Handler接口
type Handler interface {

	// User 用户相关接口
	Users() UserHandler

	// Post 帖子相关接口
	Posts() PostHandler
}

// handler 定义了Handler接口的实现
type handler struct {
	logger       *log.Logger
	store        store.Factory
	UserHandler  UserHandler
	PostHandler  PostHandler
	authStrategy middleware.AuthStrategy
}

// NewHandler 创建Handler实例
func NewHandler(logger *log.Logger, store store.Factory) Handler {

	userbiz := biz.NewUserBiz(logger, store.User())
	postbiz := biz.NewPostBiz(logger, store.Post())

	// 创建authStratrgy
	authStrategy := middleware.NewJWTStrategy(userbiz)

	h := &handler{
		logger:       logger,
		store:        store,
		authStrategy: authStrategy,
	}
	h.UserHandler = NewUserHandler(logger, userbiz, authStrategy)
	h.PostHandler = NewPostHandler(logger, postbiz)
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
