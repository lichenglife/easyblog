package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// UserHandler 用户相关接口
type UserHandler interface {

	// CreteUser 创建用户
	CreateUser() gin.HandlerFunc
	// ChangePassword 修改密码
	// ResetPassword 重置密码
	ResetPassword() gin.HandlerFunc
	// UserInfo 获取用户信息
	GetUserInfo() gin.HandlerFunc

	// UserLogin 用户登录
	UserLogin() gin.HandlerFunc
	// UserLogout 用户登出
	UserLogout() gin.HandlerFunc
	// UserInfo 获取用户信息
	UserInfo() gin.HandlerFunc
	// ListUsers 获取用户列表
	ListUsers() gin.HandlerFunc
	// GetUserByID 根据 ID 获取用户
	GetUserByID() gin.HandlerFunc
	// UpdateUser 更新用户
	UpdateUser() gin.HandlerFunc
	// DeleteUser 删除用户
	DeleteUser() gin.HandlerFunc
}

// userHandler 实现了 UserHandler 接口
type userHandler struct {
	logger  *log.Logger
	userBiz biz.UserBiz
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(logger *log.Logger, biz biz.UserBiz) UserHandler {
	return &userHandler{
		logger:  logger,
		userBiz: biz,
	}
}

var _ UserHandler = (*userHandler)(nil)

// CreateUser implements UserHandler.
func (u *userHandler) CreateUser() gin.HandlerFunc {
	panic("unimplemented")
}

// DeleteUser implements UserHandler.
func (u *userHandler) DeleteUser() gin.HandlerFunc {
	panic("unimplemented")
}

// GetUserByID implements UserHandler.
func (u *userHandler) GetUserByID() gin.HandlerFunc {
	panic("unimplemented")
}

// GetUserInfo implements UserHandler.
func (u *userHandler) GetUserInfo() gin.HandlerFunc {
	panic("unimplemented")
}

// ListUsers implements UserHandler.
func (u *userHandler) ListUsers() gin.HandlerFunc {
	panic("unimplemented")
}

// ResetPassword implements UserHandler.
func (u *userHandler) ResetPassword() gin.HandlerFunc {
	panic("unimplemented")
}

// UpdateUser implements UserHandler.
func (u *userHandler) UpdateUser() gin.HandlerFunc {
	panic("unimplemented")
}

// UserInfo implements UserHandler.
func (u *userHandler) UserInfo() gin.HandlerFunc {
	panic("unimplemented")
}

// UserLogin implements UserHandler.
func (u *userHandler) UserLogin() gin.HandlerFunc {
	panic("unimplemented")
}

// UserLogout implements UserHandler.
func (u *userHandler) UserLogout() gin.HandlerFunc {
	panic("unimplemented")
}
