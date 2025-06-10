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

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建一个新的用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param body body model.CreateUserRequest true "创建用户请求参数"
// @Success 200 {object} core.Response{data=model.UserInfo} "成功响应"
// @Failure 400 {object} core.Response "请求参数错误"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /users [post]
func (u *userHandler) CreateUser(c *gin.Context) {
	// 解析请求参数
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 检查是否为校验错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			log.Log.Error("CreateUser: validation failed", zap.Any("errors", validation.ValidationErrors(validationErrors)))
			errors := validation.ValidationErrors(validationErrors)
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

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 根据用户 ID 删除用户
// @Param userID path string true "用户 ID"
// @Tags 用户管理
// @Security     Bearer
// @Success 200 {object} core.Response "成功响应"
// @Failure 401 {object} core.Response "未授权"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Security Authorization
// @Router /users/{userID} [delete]
func (u *userHandler) DeleteUser(c *gin.Context) {
	// TODO 判断当前用户是否admin
	currentUser := c.GetString("username")
	if currentUser != "root" {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}
	userID := c.Param("userID")
	err := u.biz.UserV1().DeleteUser(c, userID)
	if err != nil {
		log.Log.Error("删除用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, "删除成功")

}

// GetUserByID 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户 ID 获取用户信息
// @Tags 用户管理
// @Security     Bearer
// @Param id path string true "用户 ID"
// @Success 200 {object} core.Response{data=model.UserInfo} "成功响应"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /users/{id} [get]
func (u *userHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("userID")
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

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Security     Bearer
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} core.Response{data=[]model.UserInfo} "成功响应"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /users [get]
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

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 根据用户 ID 重置密码
// @Tags 用户管理
// @Security Bearer
// @Param userID path string true "用户 ID"
// @Param body body model.ChangePasswordRequest true "重置密码请求参数"
// @Success 200 {object} core.Response "成功响应"
// @Failure 401 {object} core.Response "未授权"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /users/password/{userID} [put]
func (u *userHandler) ResetPassword(c *gin.Context) {
	// 1、解析参数
	userID := c.Param("userID")
	var req model.ChangePasswordRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}
	currentUser := c.GetString("userID")
	if currentUser != userID {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}
	// 2、请求biz进行密码更新
	err := u.biz.UserV1().ChangePassword(c, userID, req)
	if err != nil {
		log.Log.Error("更新用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, errno.OK, nil)

}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 根据用户 ID 更新用户信息
// @Tags 用户管理
// @Security     Bearer
// @Param userID path string true "用户 ID"
// @Param body body model.UpdateUser true "更新用户请求参数"
// @Success 200 {object} core.Response "成功响应"
// @Failure 401 {object} core.Response "未授权"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /users/{userID} [put]
func (u *userHandler) UpdateUser(c *gin.Context) {
	//1、解析参数
	userID := c.Param("userID")
	currentUser := c.GetString("userID")
	if currentUser != userID {
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

// UserLogin 用户登录
// @Summary 用户登录
// @Description 用户登录并获取 Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param body body model.UserLoginRequest true "用户登录请求参数"
// @Success 200 {object} core.Response{data=model.UserLoginResponse} "成功响应"
// @Failure 400 {object} core.Response "请求参数错误"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /users/login [post]
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

// UserLogout 用户登出
// @Summary 用户登出
// @Description 用户登出并清除 Token
// @Tags 用户管理
// @Security     Bearer
// @Success 200 {object} core.Response "成功响应"
// @Failure 500 {object} core.Response "服务器内部错误"
// @Router /users/logout [post]
func (u *userHandler) UserLogout(c *gin.Context) {
	// TODO 清除token
	c.Set("userID", "") // 清除用户ID
	c.Set("username", "")
	core.WriteResponse(c, nil, "登出成功")
}
