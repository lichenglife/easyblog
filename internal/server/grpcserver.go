package server

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	handler "github.com/lichenglife/easyblog/internal/apiserver/handler/grpc"
	"github.com/lichenglife/easyblog/internal/app"
	middleware "github.com/lichenglife/easyblog/internal/pkg/middleware/grpc"
	pb "github.com/lichenglife/easyblog/pkg/api/apiserver/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// GRPCServer 表示grpc服务器

type GRPCServer struct {
	// 配置文件
	config *viper.Viper
	// app 应用管理
	app app.IApp
	// grpcserver
	server *grpc.Server
	// 注册服务
	easyBlogService *handler.EasyBlogService
}

// NewGrpcServer
func NewGrpcServer(config *viper.Viper, app app.IApp) (*GRPCServer, error) {

	server := &GRPCServer{
		config: config,
		app:    app,
	}
	return server, nil
}

// 初始化grpc 服务器
func (s *GRPCServer) Init() error {
	// 1、认证拦截器

	// 获取认证策略
	//authStrategy := s.app.GetAuthStrategy()

	// 创建grpc 认证拦截器

	// 2、创建 grpc 服务器选项
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge:      s.config.GetDuration("server.grpc.maxConnectionAge"),
			MaxConnectionAgeGrace: s.config.GetDuration("server.grpc.maxConnectionAgeGrace"),
			MaxConnectionIdle:     s.config.GetDuration("server.grpc.maxConnectionIdle"),
			Time:                  s.config.GetDuration("server.grpc.keepAliveTime"),
			Timeout:               s.config.GetDuration("server.grpc.keepAliveTimeOut"),
		}),
		grpc.MaxRecvMsgSize(s.config.GetInt("server.grpc.maxRecvMsgSize")),
		// 开发环境不使用安全证书
		grpc.Creds(insecure.NewCredentials()),
		// 生产环境使用安全证书
		//grpc.Creds(credentials.NewTLS(&tls.Config{})),

		grpc.ChainUnaryInterceptor(
			// 请求ID 拦截器
			middleware.RequestIDInterceptor(),
			// 认证拦截器
			selector.UnaryServerInterceptor(middleware.AuthnInterceptor(), NewAuthnWhiteListMatcher()),
			// 鉴权拦截器

			// 请求默认值拦截器

			// 数据校验拦截器

		),
	}
	// 3、创建grpc server
	s.server = grpc.NewServer(opts...)

	// 4、创建处理器
	logger := s.app.GetLogger()
	store := s.app.GetStoreFactory()

	grpcHandler := handler.NewHandler(logger, store)

	//5、创建并初始化业务服务
	easyBlogService := handler.NewEasyBlogService(grpcHandler, logger)
	s.easyBlogService = easyBlogService

	//6、注册服务
	pb.RegisterEasyblogServer(s.server, easyBlogService)

	// 注册反射服务，方便使用 grpcurl 等工具调试
	reflection.Register(s.server)
	// 创建 gRPC Server → 实例化业务服务 → 注册服务 → 注册反射 → 启动监听。
	return nil
}

// NewAuthnWhiteListMatcher 创建认证白名单匹配器.
func NewAuthnWhiteListMatcher() selector.Matcher {
	whitelist := map[string]struct{}{
		pb.Easyblog_Healthz_FullMethodName:    {},
		pb.Easyblog_CreateUser_FullMethodName: {},
		pb.Easyblog_Login_FullMethodName:      {},
	}
	return selector.MatchFunc(func(ctx context.Context, call interceptors.CallMeta) bool {
		_, ok := whitelist[call.FullMethod()]
		return !ok
	})
}

// 启动服务
func (s *GRPCServer) Start() error {
	// 获取服务器端口
	port := s.config.GetInt("server.grpc.port")
	addr := fmt.Sprintf(":%d", port)

	// 监听指定端口
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("无法监听端口 %s: %v", addr, err)
	}

	// 启动服务器
	s.app.GetLogger().Info("gRPC 服务器启动成功", zap.String("addr", addr))

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.app.GetLogger().Error("gRPC 服务器运行失败", zap.Error(err))
		}
	}()

	return nil
}

// 停止服务
func (s *GRPCServer) Stop(ctx context.Context) error {
	s.app.GetLogger().Info("正在关闭 gRPC 服务器...")
	s.server.GracefulStop()
	return nil
}
