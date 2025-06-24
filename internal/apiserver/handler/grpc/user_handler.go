package grpc

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/lichenglife/easyblog/internal/pkg/middleware"
	pb "github.com/lichenglife/easyblog/pkg/api/apiserver/v1"
	"go.uber.org/zap"
)

type UserHandler interface {
	// Login 用户登录
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	// RefreshToken 刷新令牌
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error)
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)
	// CreateUser 创建用户
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	// GetUser 获取用户信息
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error)
	// ListUser 列出所有用户
	ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error)
}

// userHandler 接口实现
type userHandler struct {
	biz biz.IBiz

	authStrategy middleware.AuthStrategy
}

// ChangePassword implements UserHandler.
func (u *userHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	log.Log.Info("修改用户密码", zap.String("uesrID", req.UserID))
	//
	// 查询用户判断是否存在

	// 修改用户密码
	cpr := model.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	// 调用业务逻辑
	err := u.biz.UserV1().ChangePassword(ctx, req.UserID, cpr)
	if err != nil {
		log.Log.Error("修改用户密码失败", zap.Error(err))
		return nil, err
	}

	return &pb.ChangePasswordResponse{}, nil
}

// CreateUser implements UserHandler.
func (u *userHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Log.Info("创建用户", zap.String("username", req.Username))

	// 参数校验

	// 创建用户
	userCreate := &model.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Email:    req.Email,
	}
	user, err := u.biz.UserV1().CreateUser(ctx, userCreate)
	if err != nil {
		log.Log.Error("创建用户失败", zap.String("username", req.Username), zap.Error(err))
		return nil, err
	}

	return &pb.CreateUserResponse{UserID: user.UserID}, nil
}

// DeleteUser implements UserHandler.
func (u *userHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	panic("unimplemented")
}

// GetUser implements UserHandler.
func (u *userHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	panic("unimplemented")
}

// ListUser implements UserHandler.
func (u *userHandler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	panic("unimplemented")
}

// Login implements UserHandler.
func (u *userHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	panic("unimplemented")
}

// RefreshToken implements UserHandler.
func (u *userHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	panic("unimplemented")
}

// UpdateUser implements UserHandler.
func (u *userHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	panic("unimplemented")
}

// NewUserHandler 创建UserHandler实例

func NewUserHandler(biz biz.IBiz) UserHandler {
	authStrategy := middleware.NewJWTStrategy(biz.UserV1())
	return &userHandler{
		biz:          biz,
		authStrategy: authStrategy,
	}
}
