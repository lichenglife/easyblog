package biz

import (
	"context"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/lichenglife/easyblog/internal/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var _ UserBiz = (*userBiz)(nil)

// userBiz 定义了用户业务逻辑层
type userBiz struct {
	logger *log.Logger
	store  store.UserStore
}

// Login implements UserBiz.
func (u *userBiz) Login(ctx context.Context, username string, password string) (*model.UserInfo, error) {

	// 校验用户密码
	user, err := u.store.GetByUsername(ctx, username)
	if err != nil {
		u.logger.Logger.Error("用户登录失败", zap.Error(err))
		return nil, errno.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

		return nil, errno.ErrPasswordIncorrect
	}

	userinfo := &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}

	return userinfo, nil

}

// CreteUser implements UserBiz.
func (u *userBiz) CreteUser(ctx context.Context, user *model.User) (*model.User, error) {

	// 根据用户名称查询用户是否存在
	_, err := u.store.GetByUsername(ctx, user.Username)
	if err == nil {
		u.logger.Info("user already exist", zap.String("username", user.Username))
		return nil, errno.ErrUserAlreadyExist
	}
	if !errno.IsRecordNotFound(err) {
		return nil, err
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	// 生成用户ID
	user.UserID = utils.GenUserID()

	if err := u.store.Create(ctx, user); err != nil {
		u.logger.Error("Create user failed: %v", zap.Error(err))
		return nil, err
	}

	return user, nil
}

// DeleteUser implements UserBiz.
func (u *userBiz) DeleteUser(ctx context.Context, username string) error {
	err := u.store.Delete(ctx, username)
	if err != nil {
		u.logger.Logger.Error("删除用户失败", zap.Error(err))
		return err
	}
	return nil
}

// GetUserByID implements UserBiz.
func (u *userBiz) GetUserByID(ctx context.Context, userID string) (*model.UserInfo, error) {
	user, err := u.store.GetByID(ctx, userID)
	if err != nil {
		u.logger.Logger.Error("用户不存在", zap.String("userID", userID), zap.Error(err))
		return nil, err
	}
	return &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Nickname: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil

}

// GetUserByUsername implements UserBiz.
func (u *userBiz) GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error) {

	user, err := u.store.GetByUsername(ctx, username)
	if err != nil {
		u.logger.Error("GetUserByUsername failed: %v", zap.Error(err))
		return nil, err
	}
	return &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Nickname: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

// ListUsers implements UserBiz.
func (u *userBiz) ListUsers(ctx context.Context, page int, pageSize int) (*model.ListUserResponse, error) {
	// 请求store获取用户列表
	userlist, count, err := u.store.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	//
	users := make([]model.UserInfo, 0, len(userlist))
	for _, user := range userlist {

		users = append(users, model.UserInfo{
			UserID:    user.UserID,
			Username:  user.Username,
			Nickname:  user.NickName,
			Email:     user.Email,
			Phone:     user.Phone,
			UpdatedAt: user.UpdatedAt,
			CreatedAt: user.CreatedAt,
		})
	}
	return &model.ListUserResponse{
		Users:      users,
		TotalCount: count,
	}, nil

}

// UpdateUser implements UserBiz.
func (u *userBiz) UpdateUser(ctx context.Context, username string, req *model.UpdateUser) error {

	user, err := u.store.GetByUsername(ctx, username)
	if err != nil {
		u.logger.Error("用户不存在", zap.String("用户名", username), zap.Error(err))
		return errno.ErrNotFound
	}
	if user.Email != req.Email {
		user.Email = req.Email
	}
	if user.NickName != req.Nickname {
		user.NickName = req.Nickname
	}
	if user.Phone != req.Phone {
		user.Phone = req.Phone
	}
	user.UpdatedAt = time.Now()

	err = u.store.Update(ctx, user)
	if err != nil {
		u.logger.Logger.Error("更新用户失败", zap.String("username", username), zap.Error(err))
		return err
	}
	return nil
}

// ResetUserPassword 修改用户密码
func (u *userBiz) ResetUserPassword(ctx context.Context, username, oldPassword, newPassword string) error {
	// 1. 根据用户名获取用户
	user, err := u.store.GetByUsername(ctx, username)
	if err != nil && errno.IsRecordNotFound(err) {
		u.logger.Info("用户不存在", zap.String("username", username))
		return err
	}
	if user.Password != oldPassword {
		u.logger.Info("用户原密码错误,不允许更新")
		return errno.ErrOldPasswordIncorrect
	}

	updateUser := &model.User{
		Username: username,
		Password: newPassword,
	}
	if err := u.store.Update(ctx, updateUser); err != nil {
		u.logger.Error("更新密码失败", zap.Error(err))
	}
	return nil
}
