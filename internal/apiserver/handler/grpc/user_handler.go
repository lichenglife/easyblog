package grpc

import (
	"context"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	middleware "github.com/lichenglife/easyblog/internal/pkg/middleware/http"
	pb "github.com/lichenglife/easyblog/pkg/api/apiserver/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	log.Log.Info("删除用户请求", zap.String("userID", req.UserID))
	// 判断用户角色
	username := ctx.Value("username")
	if username != "root" {
		return nil, errno.ErrUnauthorized
	}
	// 调用biz执行删除
	err := u.biz.UserV1().DeleteUser(ctx, req.UserID)
	if err != nil {
		log.Log.Error("删除用户失败", zap.String("userID", req.UserID), zap.Error(err))
		return nil, err
	}
	return &pb.DeleteUserResponse{}, nil

}

// GetUser implements UserHandler.
func (u *userHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// get userByuserID
	userinfo, err := u.biz.UserV1().GetUserByID(ctx, req.UserID)
	if err != nil {
		log.Log.Error("删除用户失败", zap.String("userID", req.UserID), zap.Error(err))
		return nil, err
	}
	user := &pb.User{
		UserID:    userinfo.UserID,
		Username:  userinfo.Username,
		Nickname:  userinfo.Nickname,
		Email:     userinfo.Email,
		Phone:     userinfo.Phone,
		UpdatedAt: timestamppb.New(userinfo.UpdatedAt),
		CreatedAt: timestamppb.New(userinfo.CreatedAt),
	}
	return &pb.GetUserResponse{User: user}, nil

}

// ListUser implements UserHandler.
func (u *userHandler) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	if req.Offset == 0 {
		req.Offset = 1
	}
	userList, err := u.biz.UserV1().ListUsers(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		log.Log.Error("查询用户失败", zap.Error(err))
		return nil, err
	}
	users := make([]*pb.User, 0, len(userList.User))
	for _, userinfo := range userList.User {
		user := &pb.User{
			UserID:    userinfo.UserID,
			Username:  userinfo.Username,
			Nickname:  userinfo.Nickname,
			Email:     userinfo.Email,
			Phone:     userinfo.Phone,
			UpdatedAt: timestamppb.New(userinfo.UpdatedAt),
			CreatedAt: timestamppb.New(userinfo.CreatedAt),
		}
		users = append(users, user)
	}
	return &pb.ListUserResponse{
		Users:      users,
		TotalCount: userList.TotalCount,
	}, nil
}

// Login implements UserHandler.
func (u *userHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 参数校验
	// 用户校验
	ulq := model.UserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	userinfo, err := u.biz.UserV1().UserLogin(ctx, ulq)
	if err != nil {
		log.Log.Error("用户登录失败", zap.String("username", req.Username), zap.Error(err))
		return nil, err
	}
	// 生成token
	tokenString, err := u.authStrategy.GenerateToken(userinfo.UserID, userinfo.Username)
	if err != nil {
		log.Log.Error("用户登录失败", zap.String("username", req.Username), zap.Error(err))
		return nil, err
	}
	//返回数据
	return &pb.LoginResponse{
		Token:    tokenString,
		ExpireAt: timestamppb.New(time.Now().Add(time.Hour * 2)),
	}, nil
}

// RefreshToken implements UserHandler.
func (u *userHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	panic("unimplemented")
}

// UpdateUser implements UserHandler.
func (u *userHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	updateUser := &model.UpdateUser{
		UserID: req.UserID,
		Email:  req.Email,
	}
	err := u.biz.UserV1().UpdateUser(ctx, updateUser)
	if err != nil {
		log.Log.Error("更新用户失败", zap.String("userID", req.UserID), zap.Error(err))
		return nil, err
	}
	return &pb.UpdateUserResponse{}, nil

}

// NewUserHandler 创建UserHandler实例

func NewUserHandler(biz biz.IBiz) UserHandler {
	authStrategy := middleware.NewJWTStrategy(biz.UserV1())
	return &userHandler{
		biz:          biz,
		authStrategy: authStrategy,
	}
}
