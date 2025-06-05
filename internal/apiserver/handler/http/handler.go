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
	store       store.IStore
	UserHandler UserHandler
	PostHandler PostHandler
}

// NewHandler 创建Handler实例
func NewHandler(logger *log.Logger, store store.IStore) Handler {
	h := &handler{
		logger: logger,
		store:  store,
	}
	biz := biz.NewBiz(store)

	h.UserHandler = NewUserHandler(logger, biz)
	h.PostHandler = NewPostHandler(logger, biz)
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
