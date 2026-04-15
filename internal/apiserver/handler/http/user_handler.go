package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// UserHandler 用户相关接口
type UserHandler interface {
	// CreateUser 创建用户
	CreateUser(c *gin.Context)
	// UserLogin 用户登录
	UserLogin(c *gin.Context)
	// UserLogout 用户登出
	UserLogout(c *gin.Context)
	// UserInfo 获取用户信息
	UserInfo(c *gin.Context)
	// GetUserInfo 获取当前用户信息
	GetUserInfo(c *gin.Context)
	// UpdateProfile 更新用户资料（含头像）
	UpdateProfile(c *gin.Context)
	// ChangePassword 修改密码
	ChangePassword(c *gin.Context)
	// ListUsers 获取用户列表
	ListUsers(c *gin.Context)
	// GetUserByID 根据 ID 获取用户
	GetUserByID(c *gin.Context)
	// UpdateUser 更新用户
	UpdateUser(c *gin.Context)
	// DeleteUser 删除用户
	DeleteUser(c *gin.Context)
	// CreateAdmin 创建初始管理员
	CreateAdmin(c *gin.Context)
	// AuditUser 审核用户（管理员）
	AuditUser(c *gin.Context)
	// BanUser 封禁用户（管理员）
	BanUser(c *gin.Context)
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

// CreateUser 用户注册
func (u *userHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("用户注册参数验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    errno.ErrInvalidParams.Code(),
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	userInfo, err := u.userBiz.UserV1().Register(c.Request.Context(), &req)
	if err != nil {
		u.logger.Error("用户注册失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    errno.OK.Code(),
		"message": errno.OK.Message(),
		"data": gin.H{
			"user_id": userInfo.UserID,
		},
	})
}

// UserLogin 用户登录
func (u *userHandler) UserLogin(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("用户登录参数验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10003,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	resp, err := u.userBiz.UserV1().Login(c.Request.Context(), &req)
	if err != nil {
		u.logger.Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"code":    20003,
			"message": "密码错误",
			"data":    nil,
		})
		return
	}

	u.logger.Info("登录成功", zap.Any("resp", resp))
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "OK",
		"data": gin.H{
			"token": resp.Token,
			"user":  resp.User,
		},
	})
}

// UserLogout 用户登出
func (u *userHandler) UserLogout(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	// TODO: 实现 token 黑名单机制（需要 Redis 支持）
	// 将 token 加入黑名单，设置过期时间
	u.logger.Info("用户登出", zap.Any("userID", userID))

	core.WriteResponse(c, errno.OK, nil)
}

// UserInfo 获取当前登录用户信息
func (u *userHandler) UserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	userInfo, err := u.userBiz.UserV1().GetUserByUsername(c.Request.Context(), username.(string))
	if err != nil {
		u.logger.Error("获取用户信息失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, userInfo)
}

// GetUserInfo 获取当前用户信息（通过用户 ID）
func (u *userHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	userInfo, err := u.userBiz.UserV1().GetUserByUUID(c.Request.Context(), userID.(string))
	if err != nil {
		u.logger.Error("获取用户信息失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	u.logger.Info("获取用户信息成功", zap.Any("userInfo", userInfo))
	c.JSON(http.StatusOK, gin.H{
		"code":    errno.OK.Code(),
		"message": errno.OK.Message(),
		"data":    userInfo,
	})
}

// UpdateProfile 更新用户资料（含头像）
func (u *userHandler) UpdateProfile(c *gin.Context) {
	// 从上下文获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	var req struct {
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
		Bio      string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("更新资料参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := u.userBiz.UserV1().UpdateProfile(c.Request.Context(), userID.(string), req.Avatar, req.Nickname, req.Bio); err != nil {
		u.logger.Error("更新用户资料失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// ChangePassword 修改密码
func (u *userHandler) ChangePassword(c *gin.Context) {
	// 从上下文获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		core.WriteResponse(c, errno.ErrUnauthorized, nil)
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("修改密码参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := u.userBiz.UserV1().ChangePassword(c.Request.Context(), userID.(string), req.OldPassword, req.NewPassword); err != nil {
		u.logger.Error("修改密码失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// ListUsers 获取用户列表
func (u *userHandler) ListUsers(c *gin.Context) {
	page := 1
	pageSize := 10

	// 解析页码参数
	if p := c.Query("page"); p != "" {
		if _, err := fmt.Sscanf(p, "%d", &page); err != nil {
			u.logger.Warn("页码参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil {
			u.logger.Warn("页面大小参数格式错误", zap.Error(err))
			core.WriteResponse(c, errno.ErrInvalidParams, nil)
			return
		}
	}

	// 限制最大页面大小
	if pageSize > 100 {
		pageSize = 100
	}

	users, err := u.userBiz.UserV1().ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		u.logger.Error("获取用户列表失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, users)
}

// GetUserByID 根据 ID 获取用户
func (u *userHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		u.logger.Warn("用户 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	// 将字符串 userID 转换为 uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		u.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	userInfo, err := u.userBiz.UserV1().GetUserByID(c.Request.Context(), id)
	if err != nil {
		u.logger.Error("获取用户信息失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, userInfo)
}

// UpdateUser 更新用户
func (u *userHandler) UpdateUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		u.logger.Warn("用户名不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var req model.UpdateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("更新用户参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := u.userBiz.UserV1().UpdateUser(c.Request.Context(), username, &req); err != nil {
		u.logger.Error("更新用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// DeleteUser 删除用户
func (u *userHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		u.logger.Warn("用户 ID 不能为空")
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		u.logger.Warn("用户 ID 格式错误", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := u.userBiz.UserV1().DeleteUser(c.Request.Context(), id); err != nil {
		u.logger.Error("删除用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// AuditUser 审核用户（管理员）
func (u *userHandler) AuditUser(c *gin.Context) {
	var req model.AuditUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("审核用户参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := u.userBiz.UserV1().AuditUser(c.Request.Context(), req.UserID, req.Status); err != nil {
		u.logger.Error("审核用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// BanUser 封禁用户（管理员）
func (u *userHandler) BanUser(c *gin.Context) {
	var req model.BanUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("封禁用户参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	if err := u.userBiz.UserV1().BanUser(c.Request.Context(), req.UserID); err != nil {
		u.logger.Error("封禁用户失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, nil)
}

// CreateAdmin 创建初始管理员
func (u *userHandler) CreateAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Password string `json:"password" binding:"required,min=6,max=30"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		u.logger.Warn("创建管理员参数验证失败", zap.Error(err))
		core.WriteResponse(c, errno.ErrInvalidParams, nil)
		return
	}

	userInfo, err := u.userBiz.UserV1().CreateAdmin(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		u.logger.Error("创建管理员失败", zap.Error(err))
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, errno.OK, userInfo)
}
