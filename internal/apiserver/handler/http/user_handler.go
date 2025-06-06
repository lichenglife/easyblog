package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// UserHandler 用户相关接口
type UserHandler interface {

	// CreteUser 创建用户
	CreateUser(c *gin.Context)
	// ChangePassword 修改密码
	// ResetPassword 重置密码
	ResetPassword(c *gin.Context)
	// UserInfo 获取用户信息
	GetUserInfo(c *gin.Context)

	// UserLogin 用户登录
	UserLogin(c *gin.Context)
	// UserLogout 用户登出
	UserLogout(c *gin.Context)
	// UserInfo 获取用户信息
	UserInfo(c *gin.Context)
	// ListUsers 获取用户列表
	ListUsers(c *gin.Context)
	// GetUserByID 根据 ID 获取用户
	GetUserByID(c *gin.Context)
	// UpdateUser 更新用户
	UpdateUser(c *gin.Context)
	// DeleteUser 删除用户
	DeleteUser(c *gin.Context)
}

// userHandler 实现了 UserHandler 接口
type userHandler struct {
	logger  *log.Logger
	userBiz biz.IBiz
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(logger *log.Logger, biz biz.IBiz) UserHandler {
	return &userHandler{
		logger:  logger,
		userBiz: biz,
	}
}

var _ UserHandler = (*userHandler)(nil)

// CreateUser implements UserHandler.
func (u *userHandler) CreateUser(c *gin.Context) {
	// 解析请求参数
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

// DeleteUser implements UserHandler.
func (u *userHandler) DeleteUser(c *gin.Context) {
	panic("unimplemented")
}

// GetUserByID implements UserHandler.
func (u *userHandler) GetUserByID(c *gin.Context) {
	panic("unimplemented")
}

// GetUserInfo implements UserHandler.
func (u *userHandler) GetUserInfo(c *gin.Context) {
	panic("unimplemented")
}

// ListUsers implements UserHandler.
func (u *userHandler) ListUsers(c *gin.Context) {
	panic("unimplemented")
}

// ResetPassword implements UserHandler.
func (u *userHandler) ResetPassword(c *gin.Context) {
	panic("unimplemented")
}

// UpdateUser implements UserHandler.
func (u *userHandler) UpdateUser(c *gin.Context) {
	panic("unimplemented")
}

// UserInfo implements UserHandler.
func (u *userHandler) UserInfo(c *gin.Context) {
	panic("unimplemented")
}

// UserLogin implements UserHandler.
func (u *userHandler) UserLogin(c *gin.Context) {
	panic("unimplemented")
}

// UserLogout implements UserHandler.
func (u *userHandler) UserLogout(c *gin.Context) {
	panic("unimplemented")
}
