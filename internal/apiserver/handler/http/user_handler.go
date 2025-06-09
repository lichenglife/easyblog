package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/lichenglife/easyblog/internal/pkg/middleware"
	"github.com/lichenglife/easyblog/internal/pkg/validation"
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
	logger *log.Logger
	biz    biz.IBiz

	authStrategy middleware.AuthStrategy
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(logger *log.Logger, biz biz.IBiz) UserHandler {
	authstraty := middleware.NewJWTStrategy(biz.UserV1())
	return &userHandler{
		logger:       logger,
		biz:          biz,
		authStrategy: authstraty,
	}
}

var _ UserHandler = (*userHandler)(nil)

// CreateUser implements UserHandler.
func (u *userHandler) CreateUser(c *gin.Context) {
	// 解析请求参数
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 检查是否为校验错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			log.Log.Error("CreateUser: validation failed", zap.Any("errors", validation.ValidationErrors(validationErrors)))
			var errors map[string]string = make(map[string]string)
			errors = validation.ValidationErrors(validationErrors)
			core.WriteResponse(c, nil, errors)
			// core.WriteResponse(c, gin.H{
			// 	"code":    10001,
			// 	"message": "参数校验失败",
			// 	"errors":  validation.ValidationErrors(validationErrors),
			// }, nil)
			return
		}

		// 其他错误
		log.Log.Error("CreateUser: failed to bind JSON", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	// 创建用户
	var user *model.UserInfo
	user, err := u.biz.UserV1().CreateUser(c, &req)
	if err != nil {
		log.Log.Error("CreateUser: failed to create user", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, user)
}

// DeleteUser implements UserHandler.
func (u *userHandler) DeleteUser(c *gin.Context) {
	// TODO 判断当前用户是否admin
	currentUser, ok := c.Get("username")
	if !ok || currentUser != "root" {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}
	userID := c.Param("id")
	err := u.biz.UserV1().DeleteUser(c, userID)
	if err != nil {
		log.Log.Error("删除用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, "删除成功")

}

// GetUserByID implements UserHandler.
func (u *userHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	userInfo, err := u.biz.UserV1().GetUserByID(c, userID)
	if err != nil {
		log.Log.Error("获取用户信息失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, userInfo)
}

// GetUserInfo implements UserHandler.
func (u *userHandler) GetUserInfo(c *gin.Context) {
	panic("unimplemented")
}

// ListUsers implements UserHandler.
func (u *userHandler) ListUsers(c *gin.Context) {
	// 获取分页查询参数
	// 1. 获取当前页码
	limit := core.GetLimitParam(c)
	page := core.GetPageParam(c)
	// 从当前c 中获取 limit 、page
	userList, err := u.biz.UserV1().ListUsers(c, page, limit)
	if err != nil {
		log.Log.Error("获取用户列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, userList)

}

// ResetPassword implements UserHandler.
func (u *userHandler) ResetPassword(c *gin.Context) {
	// 1、解析参数
	userID := c.Param("userID")
	var req model.ChangePasswordRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}
	currentUser, ok := c.Get("userID")
	if !ok || currentUser != userID {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}
	// 2、请求biz进行密码跟新
	err := u.biz.UserV1().ChangePassword(c, userID, req)
	if err != nil {
		log.Log.Error("更新用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, errno.OK, nil)

}

// UpdateUser implements UserHandler.
func (u *userHandler) UpdateUser(c *gin.Context) {
	//1、解析参数
	userID := c.Param("userID")
	currentUser, ok := c.Get("userID")
	if !ok || currentUser != userID {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	var updateUserModel model.UpdateUser
	if err := c.ShouldBindJSON(&updateUserModel); err != nil {
		errmessages := validation.ValidationErrors(err)
		core.WriteResponse(c, nil, errmessages)
		return
	}
	updateUserModel.UserID = userID
	// 2、更新用户
	err := u.biz.UserV1().UpdateUser(c, &updateUserModel)
	if err != nil {
		log.Log.Error("更新用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, "更新成功")
}

// UserInfo implements UserHandler.
func (u *userHandler) UserInfo(c *gin.Context) {
	panic("unimplemented")
}

// UserLogin implements UserHandler.
func (u *userHandler) UserLogin(c *gin.Context) {
	var userLogin model.UserLoginRequest

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		errmessages := validation.ValidationErrors(err)
		core.WriteResponse(c, errno.ErrInvalidParams, errmessages)
		return
	}
	userinfo, err := u.biz.UserV1().UserLogin(c, userLogin)
	if err != nil {
		log.Log.Error("用户登录失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	// 生成token
	tokenString, err := u.authStrategy.GenerateToken(userinfo.UserID, userinfo.Username)
	if err != nil {
		log.Log.Error("生成token失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrGenerateToken, nil)
		return
	}
	userLoginResponse := &model.UserLoginResponse{
		User:  *userinfo,
		Token: tokenString,
	}

	core.WriteResponse(c, nil, userLoginResponse)

}

// UserLogout implements UserHandler.
func (u *userHandler) UserLogout(c *gin.Context) {
	panic("unimplemented")
}
