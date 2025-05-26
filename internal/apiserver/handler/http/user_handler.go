package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/lichenglife/easyblog/internal/pkg/middleware"
	"go.uber.org/zap"
)

// UserHandler 用户相关接口
type UserHandler interface {

	// CreteUser 创建用户
	CreateUser(c *gin.Context)
	// ChangePassword 修改密码
	// ResetPassword 重置密码
	ResetPassword(c *gin.Context)
	// UserInfo 获取用户信息
	// GetUserInfo(c *gin.Context)

	// UserLogin 用户登录
	UserLogin(c *gin.Context)
	// UserLogout 用户登出
	UserLogout(c *gin.Context)

	// ListUsers 获取用户列表
	ListUsers(c *gin.Context)
	// GetUserByUsername 根据 用户名称 获取用户
	GetUserByUsername(c *gin.Context)
	// UpdateUser 更新用户
	UpdateUser(c *gin.Context)
	// DeleteUser 删除用户
	DeleteUser(c *gin.Context)

	GetUserInfo(c *gin.Context)
}

// userHandler 实现了 UserHandler 接口
type userHandler struct {
	logger       *log.Logger
	userBiz      biz.UserBiz
	authStrategy middleware.AuthStrategy
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(logger *log.Logger, biz biz.UserBiz, authStrategy middleware.AuthStrategy) UserHandler {
	return &userHandler{
		logger:       logger,
		userBiz:      biz,
		authStrategy: authStrategy,
	}
}

var _ UserHandler = (*userHandler)(nil)

// CreateUser implements UserHandler.
func (u *userHandler) CreateUser(c *gin.Context) {

	// 获取请求参数
	var userRequest *model.CreateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		u.logger.Error("CreateUser bind json error", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	user := &model.User{
		Username:  userRequest.Username,
		NickName:  userRequest.Nickname,
		Email:     userRequest.Email,
		Phone:     userRequest.Phone,
		Password:  userRequest.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 调用biz层逻辑
	userInfo, err := u.userBiz.CreteUser(c, user)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, userInfo)

}

// DeleteUser implements UserHandler.
func (u *userHandler) DeleteUser(c *gin.Context) {
	// 根据用户名删除用户
	// 判断当前用户是否管理员用户
	currentUser := c.GetString("username")
	if currentUser != "admin" {
		core.WriteResponse(c, errno.ErrForbidden, nil)
		u.logger.Logger.Error("删除用户失败，不是管理员用户")
		return
	}
	username := c.Query("username")

	if err := u.userBiz.DeleteUser(c, username); err != nil {

	}

}

// GetUserByID implements UserHandler.
func (u *userHandler) GetUserByUsername(c *gin.Context) {
	// 获取请求参数
	username := c.Param("username")
	// 调用biz层逻辑
	userInfo, err := u.userBiz.GetUserByUsername(c, username)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, userInfo)
}

// GetUserInfo implements UserHandler.
func (u *userHandler) GetUserInfo(c *gin.Context) {
	userid := c.Param("userid")

	user, err := u.userBiz.GetUserByID(c, userid)
	if err != nil {

		return
	}

	core.WriteResponse(c, nil, user)

}

// ListUsers implements UserHandler.
func (u *userHandler) ListUsers(c *gin.Context) {

	u.logger.Info("ListUsers", zap.Int("page", core.GetPageParam(c)), zap.Int("pagesize", core.GetLimitParam(c)))
	page := core.GetPageParam(c)
	pagesize := core.GetLimitParam(c)
	users, err := u.userBiz.ListUsers(c, page, pagesize)
	if err != nil {
		u.logger.Error("ListUsers failed", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, users)
}

// ResetPassword implements UserHandler.
func (u *userHandler) ResetPassword(c *gin.Context) {

	username := c.Query("username")

	var userRequest *model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		u.logger.Error("ResetPassword bind json error", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	// 调用biz层逻辑
	if err := u.userBiz.ResetUserPassword(c, username, userRequest.OldPassword, userRequest.NewPassword); err != nil {
		u.logger.Error("ResetPassword failed", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)
}

// UpdateUser implements UserHandler.
func (u *userHandler) UpdateUser(c *gin.Context) {
	username := c.Param("username")

	var user *model.UpdateUser

	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Logger.Error("UpdateUser bind json error", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	if err := u.userBiz.UpdateUser(c, username, user); err != nil {
		u.logger.Logger.Error("UpdateUser failed", zap.Error(err))
		core.WriteResponse(c, err, "更新用户失败")
		return
	}

	core.WriteResponse(c, nil, "更新成功")
}

// UserLogin implements UserHandler.
func (u *userHandler) UserLogin(c *gin.Context) {
	// 账号密码校验
	var user model.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Logger.Error("Binding UserLogin json failed", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	userInfo, err := u.userBiz.Login(c, user.Username, user.Password)
	if err != nil {
		u.logger.Logger.Error("用户登录认证失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	// 生成token并保存到redis中

	tokenString, err := u.authStrategy.GenerateToken(userInfo.UserID, user.Username)
	if err != nil {
		u.logger.Logger.Error("用户登录授权失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	userLogin := &model.UserLoginResponse{
		User:  *userInfo,
		Token: tokenString,
	}
	core.WriteResponse(c, nil, userLogin)

}

// UserLogout implements UserHandler.
func (u *userHandler) UserLogout(c *gin.Context) {
	// 判断登录状态

	// 删除token
}
