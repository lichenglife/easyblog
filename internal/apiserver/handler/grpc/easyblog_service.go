package grpc

import (
	"context"

	"github.com/lichenglife/easyblog/internal/pkg/log"
	pb "github.com/lichenglife/easyblog/pkg/api/apiserver/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// EasyBlogService 实现apiserver grcp 接口， grpc请求入口
/**
  1、 grpc框架自动将客户请求路由到EasyBlog的对应方法
  2、 请求转发到对应的handler
  3、再转发到具体的biz层/store层处理请求
  4、响应返回
*/
type EasyBlogService struct {
	// 服务端实现proto文件中定义的服务接口(接口嵌入)
	pb.UnimplementedEasyblogServer
	// handler 处理层
	handler Handler
	logger  *log.Logger
}

// NewEasyBlogService 实例化EasyBlogServie 实例

func NewEasyBlogService(handler Handler, logger *log.Logger) *EasyBlogService {
	return &EasyBlogService{
		handler: handler,
		logger:  logger,
	}
}

func (e *EasyBlogService) Healthz(ctx context.Context, req *emptypb.Empty) (*pb.HealthzResponse, error) {
	e.logger.Logger.Info("健康检查请求...")
	return e.handler.Healthz().Healthz(ctx, req)
}
func (e *EasyBlogService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return e.handler.Users().Login(ctx, req)
}
func (e *EasyBlogService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return e.handler.Users().RefreshToken(ctx, req)
}
func (e *EasyBlogService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return e.handler.Users().ChangePassword(ctx, req)
}
func (e *EasyBlogService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return e.handler.Users().CreateUser(ctx, req)
}
func (e *EasyBlogService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return e.handler.Users().UpdateUser(ctx, req)
}
func (e *EasyBlogService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return e.handler.Users().DeleteUser(ctx, req)
}
func (e *EasyBlogService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return e.handler.Users().GetUser(ctx, req)
}
func (e *EasyBlogService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	return e.handler.Users().ListUser(ctx, req)
}
func (e *EasyBlogService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	return e.handler.Posts().CreatePost(ctx, req)
}
func (e *EasyBlogService) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	return e.handler.Posts().UpdatePost(ctx, req)
}
func (e *EasyBlogService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	return e.handler.Posts().DeletePost(ctx, req)
}
func (e *EasyBlogService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	return e.handler.Posts().GetPost(ctx, req)
}
func (e *EasyBlogService) ListPost(ctx context.Context, req *pb.ListPostRequest) (*pb.ListPostResponse, error) {
	return e.handler.Posts().ListPost(ctx, req)
}
