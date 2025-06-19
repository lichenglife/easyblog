package server

import (
	"context"
	"fmt"

	"github.com/lichenglife/easyblog/internal/app"
	"github.com/lichenglife/easyblog/internal/pkg/monitor"
	"github.com/spf13/viper"
)

// 定义IServer 接口

type IServer interface {

	// 初始化Server
	Init() error

	// 启动服务
	Start() error

	// 停止服务
	Stop(ctx context.Context) error
}

// UnionServer 代表一个统一的server， 包括httpserver 、grpcserver
type UnionServer struct {
	//  配置
	cfg *viper.Viper
	app app.IApp
	// httpserver
	httpServer *HTTPServer
	// 监控
	monitor *monitor.Monitor
	//  grpcserver
}

func NewUnionServer(cfg *viper.Viper, app app.IApp) (*UnionServer, error) {
	server := &UnionServer{
		cfg: cfg,
		app: app,
	}

	// 初始化监控
	monitor := monitor.NewMonitor(
		"easyblog",
		cfg.GetString("server.http.addr"),
		cfg.GetBool("monitor.pprof.enable"),
		cfg.GetBool("monitor.trace.enable"),
		cfg.GetBool("monitor.metrics.enable"),
	)
	if err := monitor.Init(); err != nil {
		return nil, fmt.Errorf("初始化监控失败: %v", err)
	}
	server.monitor = monitor

	// 初始化 httpServer
	httpServer, err := NewHttpServer(cfg, app)
	if err != nil {
		return nil, fmt.Errorf("初始化httpServer失败%v", err)
	}
	server.httpServer = httpServer

	// 初始化grpcServer

	return server, nil
}

func (s *UnionServer) Init() error {
	// 初始化httpserver
	if err := s.httpServer.Init(); err != nil {
		return fmt.Errorf("初始化HTTPServer失败%v", err)
	}

	// 注册监控路由
	s.monitor.RegisterRoutes(s.httpServer.engine)

	// 添加监控中间件
	s.httpServer.engine.Use(s.monitor.Middleware())

	return nil
}

func (s *UnionServer) Start() error {

	// 启动HTTPServer 服务
	if err := s.httpServer.Start(); err != nil {
		return fmt.Errorf("启动HTTPServer失败%v", err)
	}
	// 启动GRPCServer 服务
	return nil
}

func (s *UnionServer) Stop(ctx context.Context) error {

	s.app.GetLogger().Logger.Info("正在停止服务")
	if err := s.httpServer.Stop(ctx); err != nil {

		return fmt.Errorf("停止HTTPServer失败%v", err)

	}

	return nil
}
