package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/app"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// HttpServer http服务

type HTTPServer struct {
	// 配置
	config *viper.Viper
	// app  app中封装了日志、storefactory
	app app.IApp
	// gin
	engine *gin.Engine
	// http
	http *http.Server
	// handler 处理器
}

func NewHttpServer(config *viper.Viper, app app.IApp) (*HTTPServer, error) {

	server := &HTTPServer{
		config: config,
		app:    app,
	}

	// todo 创建业务处理器handler

	return server, nil
}

// Init 初始化httpServer
func (s *HTTPServer) Init() error {
	//  初始化 engine
	err := s.initEngine()
	if err != nil {
		return fmt.Errorf("初始化Engine失败:%v", err)
	}
	// 注册http路由
	if err := s.registerRoutes(); err != nil {
		return fmt.Errorf("注册路由规则失败:%v", err)
	}

	// 初始化httpserver
	if err := s.initHttpServer(); err != nil {
		return fmt.Errorf("初始化httpserver失败:%v", err)
	}
	return nil
}

// initEngine 初始化engine实例

func (s *HTTPServer) initEngine() error {

	// 设置运行模式
	gin.SetMode(s.config.GetString("server.http.mode"))

	engine := gin.New()

	// 注册中间件
	engine.Use(
		// RequestID
		middleware.RequestID(),
		// 请求日志记录
		middleware.Logger(s.app.GetLogger()),
		// 故障恢复
		middleware.Recovery(s.app.GetLogger()),
		// 跨域
		middleware.CORS(),
	)
	s.engine = engine

	return nil
}

func (s *HTTPServer) registerRoutes() error {

	// 健康检查
	s.engine.GET("/healthz", s.healthcheck)

	// swagger api接口文档

	// 业务功能路由规则
	// v1 := s.engine.Group("/v1")
	// {
	// 	//
	// }

	return nil
}

func (s *HTTPServer) healthcheck(c *gin.Context) {

	core.WriteResponse(c, nil, gin.H{"status": "OK"})
}

// 初始化httpserver
func (s *HTTPServer) initHttpServer() error {

	s.http = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.GetString("server.http.host"), s.config.GetInt("server.http.port")),
		Handler:        s.engine,
		ReadTimeout:    time.Duration(s.config.GetInt("server.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(s.config.GetInt("server.http.write_timeout")) * time.Second,
		MaxHeaderBytes: s.config.GetInt("server.http.max_header_bytes"),
	}
	return nil

}

// Start 启动服务
func (s *HTTPServer) Start() error {

	logger := s.app.GetLogger().Logger
	logger.Info("启动HTTPServer", zap.String("addr", s.http.Addr), zap.String("mode", s.config.GetString("server.http.mode")))

	go func() {
		if err := s.http.ListenAndServe(); err != nil {
			s.app.GetLogger().Error("启动HTTP服务器失败", zap.Error(err))
		}
	}()
	return nil
}

// Stop 停止服务
func (s *HTTPServer) Stop(ctx context.Context) error {

	s.app.GetLogger().Logger.Info("正在停止HTTP服务器...")
	if err := s.http.Shutdown(ctx); err != nil {
		s.app.GetLogger().Logger.Error("HTTP服务器停止失败", zap.Error(err))
	}
	return nil
}
