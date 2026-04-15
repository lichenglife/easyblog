package user

import (
	"context"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/auth"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

// UserBiz 用户业务接口
type UserBiz interface {
	// Register 用户注册
	Register(ctx context.Context, req *model.CreateUserRequest) (*model.UserInfo, error)
	// Login 用户登录
	Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserLoginResponse, error)
	// GetByID 根据 ID 获取用户
	GetUserByID(ctx context.Context, id uint) (*model.UserInfo, error)
	// GetByUUID 根据 UUID 获取用户
	GetUserByUUID(ctx context.Context, uuid string) (*model.UserInfo, error)
	// GetByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error)
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, userID string, user *model.UpdateUser) error
	// UpdateProfile 更新用户资料（含头像）
	UpdateProfile(ctx context.Context, userID string, avatar, nickname, bio string) error
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id uint) error
	// ListUsers 获取用户列表
	ListUsers(ctx context.Context, page, pageSize int) (*model.ListUserResponse, error)
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error
	// AuditUser 审核用户
	AuditUser(ctx context.Context, userID uint, status int) error
	// BanUser 封禁用户
	BanUser(ctx context.Context, userID uint) error
	// CreateAdmin 创建初始管理员
	CreateAdmin(ctx context.Context, username, password, email string) (*model.UserInfo, error)
}

func NewUserBiz(store store.UserStore, jwt *auth.JWT) UserBiz {
	return &userBiz{
		store: store,
		jwt:   jwt,
	}
}

var _ UserBiz = (*userBiz)(nil)

// userBiz 实现了 UserBiz 接口
type userBiz struct {
	store store.UserStore
	jwt   *auth.JWT
}

// Register implements UserBiz.
func (u *userBiz) Register(ctx context.Context, req *model.CreateUserRequest) (*model.UserInfo, error) {
	// 检查用户名是否已存在
	existingUser, err := u.store.GetByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, errno.ErrUserAlreadyExist
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errno.ErrInternalServer.WithMessage("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		UserID:   util.GenerateUUID(),
		Username: req.Username,
		Password: string(hashedPassword),
		NickName: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     1,   // 普通用户
		Status:   0,   // 待审核
	}

	if err := u.store.Create(ctx, user); err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Nickname: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Status:   user.Status,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		CreateAt: user.CreateAt.Format(time.RFC3339),
	}, nil
}

// Login implements UserBiz.
func (u *userBiz) Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	// 获取用户
	user, err := u.store.GetByUsername(ctx, req.Username)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// 生成 Token
	token, err := u.jwt.CreateToken(user.UserID, user.Username, user.Role)
	if err != nil {
		return nil, errno.ErrInternalServer.WithMessage("Token 生成失败")
	}

	resp := &model.UserLoginResponse{
		Token: token,
		User: model.UserInfo{
			UserID:   user.UserID,
			Username: user.Username,
			Nickname: user.NickName,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			Status:   user.Status,
			Avatar:   user.Avatar,
			Bio:      user.Bio,
			CreateAt: user.CreateAt.Format(time.RFC3339),
		},
	}
	return resp, nil
}

// GetUserByID implements UserBiz.
func (u *userBiz) GetUserByID(ctx context.Context, id uint) (*model.UserInfo, error) {
	user, err := u.store.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Nickname: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Status:   user.Status,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		CreateAt: user.CreateAt.Format(time.RFC3339),
	}, nil
}

// GetUserByUUID 根据 UUID 获取用户
func (u *userBiz) GetUserByUUID(ctx context.Context, uuid string) (*model.UserInfo, error) {
	user, err := u.store.GetByUUID(ctx, uuid)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Nickname: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Status:   user.Status,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		CreateAt: user.CreateAt.Format(time.RFC3339),
	}, nil
}

// GetUserByUsername implements UserBiz.
func (u *userBiz) GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error) {
	user, err := u.store.GetByUsername(ctx, username)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Nickname: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Status:   user.Status,
		Avatar:   user.Avatar,
		Bio:      user.Bio,
		CreateAt: user.CreateAt.Format(time.RFC3339),
	}, nil
}

// UpdateUser implements UserBiz.
func (u *userBiz) UpdateUser(ctx context.Context, userID string, user *model.UpdateUser) error {
	// 获取用户
	existingUser, err := u.store.GetByUsername(ctx, userID)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrUserNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 更新用户信息
	existingUser.NickName = user.Nickname
	existingUser.Email = user.Email
	existingUser.Phone = user.Phone

	if err := u.store.Update(ctx, existingUser); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// DeleteUser implements UserBiz.
func (u *userBiz) DeleteUser(ctx context.Context, id uint) error {
	if err := u.store.Delete(ctx, id); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}
	return nil
}

// ListUsers implements UserBiz.
func (u *userBiz) ListUsers(ctx context.Context, page, pageSize int) (*model.ListUserResponse, error) {
	users, err := u.store.List(ctx, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	var userInfos []model.UserInfo
	for _, user := range users {
		userInfos = append(userInfos, model.UserInfo{
			UserID:   user.UserID,
			Username: user.Username,
			Nickname: user.NickName,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			Status:   user.Status,
			Avatar:   user.Avatar,
			Bio:      user.Bio,
			CreateAt: user.CreateAt.Format(time.RFC3339),
		})
	}

	return &model.ListUserResponse{
		TotalCount: int64(len(userInfos)),
		HasMore:    len(userInfos) == pageSize,
		User:       userInfos,
	}, nil
}

// ChangePassword 修改密码
func (u *userBiz) ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	// 获取用户
	user, err := u.store.GetByUUID(ctx, userID)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrUserNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errno.ErrPasswordIncorrect
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errno.ErrInternalServer.WithMessage("密码加密失败")
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := u.store.Update(ctx, user); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// AuditUser 审核用户
func (u *userBiz) AuditUser(ctx context.Context, userID uint, status int) error {
	user, err := u.store.GetByID(ctx, userID)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrUserNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	user.Status = status
	if err := u.store.Update(ctx, user); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// BanUser 封禁用户
func (u *userBiz) BanUser(ctx context.Context, userID uint) error {
	user, err := u.store.GetByID(ctx, userID)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrUserNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	user.Status = 2 // 2-封禁
	if err := u.store.Update(ctx, user); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// UpdateProfile 更新用户资料（含头像）
func (u *userBiz) UpdateProfile(ctx context.Context, userID string, avatar, nickname, bio string) error {
	user, err := u.store.GetByUUID(ctx, userID)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrUserNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	if avatar != "" {
		user.Avatar = avatar
	}
	if nickname != "" {
		user.NickName = nickname
	}
	if bio != "" {
		user.Bio = bio
	}

	if err := u.store.Update(ctx, user); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// CreateAdmin 创建初始管理员
func (u *userBiz) CreateAdmin(ctx context.Context, username, password, email string) (*model.UserInfo, error) {
	// 检查是否已存在管理员
	admins, err := u.store.List(ctx, 1, 1)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	// 简单检查：如果已有用户且用户名是 admin，则不再创建
	if len(admins) > 0 {
		existingAdmin, _ := u.store.GetByUsername(ctx, username)
		if existingAdmin != nil && existingAdmin.Role == 2 {
			return nil, errno.ErrUserAlreadyExist.WithMessage("管理员已存在")
		}
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errno.ErrInternalServer.WithMessage("密码加密失败")
	}

	// 创建管理员
	admin := &model.User{
		UserID:   util.GenerateUUID(),
		Username: username,
		Password: string(hashedPassword),
		NickName: "管理员",
		Email:    email,
		Role:     2,   // 管理员
		Status:   1,   // 正常
	}

	if err := u.store.Create(ctx, admin); err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.UserInfo{
		UserID:   admin.UserID,
		Username: admin.Username,
		Nickname: admin.NickName,
		Email:    admin.Email,
		Role:     admin.Role,
		Status:   admin.Status,
		Avatar:   admin.Avatar,
		Bio:      admin.Bio,
		CreateAt: admin.CreateAt.Format(time.RFC3339),
	}, nil
}
