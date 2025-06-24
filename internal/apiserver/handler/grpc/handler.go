package grpc

import (
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// Handler 定义处理grpc的接口

type Handler interface {

	// Healthz Handler

	Healthz() HealthzHandler

	// User Handler
	Users() UserHandler

	// Post Handler
	Posts() PostHandler
}

// Handler 的接口实现
type handler struct {
	logger        *log.Logger
	healthHandler HealthzHandler
	userHandler   UserHandler
	postHandler   PostHandler
}

// Healthz implements Handler.
func (h *handler) Healthz() HealthzHandler {
	return h.healthHandler
}

// Posts implements Handler.
func (h *handler) Posts() PostHandler {
	return h.postHandler
}

// Users implements Handler.
func (h *handler) Users() UserHandler {
	return h.userHandler
}

// NewHandler 创建Handler 接口实现
func NewHandler(log *log.Logger, store store.IStore) Handler {
	biz := biz.NewBiz(store)
	userHandler := NewUserHandler(biz)
	postHandler := NewPostHandler(biz)
	healthHandler := NewHealzHandler(log)
	return &handler{
		logger:        log,
		healthHandler: healthHandler,
		userHandler:   userHandler,
		postHandler:   postHandler,
	}
}
