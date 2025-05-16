package handler

import (
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/log"
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
	logger      *log.Logger
	store       store.Factory
	UserHandler UserHandler
	PostHandler PostHandler
}

// NewHandler 创建Handler实例
func NewHandler(logger *log.Logger, store store.Factory) Handler {
	h := &handler{
		logger: logger,
		store:  store,
	}

	userbiz := biz.NewUserBiz(logger, store.User())
	postbiz := biz.NewPostBiz(logger, store.Post())

	h.UserHandler = NewUserHandler(logger, userbiz)
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
