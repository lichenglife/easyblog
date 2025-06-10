package biz

import (
	"context"

	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/lichenglife/easyblog/internal/pkg/utils/authn"
	genid "github.com/lichenglife/easyblog/internal/pkg/utils/genID"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// UserBiz 用户业务接口
type UserBiz interface {
	// Create 创建用户
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserInfo, error)
	// GetByID 根据 ID 获取用户
	GetUserByID(ctx context.Context, userID string) (*model.UserInfo, error)
	// GetByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error)
	// Update 更新用户
	UpdateUser(ctx context.Context, user *model.UpdateUser) error
	// Delete 删除用户
	DeleteUser(ctx context.Context, userID string) error
	// List 获取用户列表
	ListUsers(ctx context.Context, page, pageSize int) (*model.ListUserResponse, error)
	// UserLogin 用户登录
	UserLogin(ctx context.Context, user model.UserLoginRequest) (*model.UserInfo, error)

	// ChangePassword 更新用户密码
	ChangePassword(ctx context.Context, userID string, user model.ChangePasswordRequest) error
}

func NewUserBiz(store store.IStore) UserBiz {
	return &userBiz{
		store: store,
	}
}

var _ UserBiz = (*userBiz)(nil)

// userBiz 定义了用户业务逻辑层
type userBiz struct {
	store store.IStore
}

// ChangePassword implements UserBiz.
func (u *userBiz) ChangePassword(ctx context.Context, userID string, req model.ChangePasswordRequest) error {
	// 1、查询用户
	user, err := u.store.User().GetByID(ctx, userID)
	if err != nil {
		return err
	}
	// 2、判断密码是否一致
	if err := authn.Compare(user.Password, req.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}
	password, err := authn.Encrypt(req.NewPassword)
	if err  != nil {
		return  err
	}
	
	// 3、更新用户密码
	updateUser := model.User{
		UserID:   userID,
		Password: password,
	}
	if err := u.store.User().Update(ctx, &updateUser); err != nil {
		log.Log.Error(err.Error())
		return err
	}
	return nil

}

// UserLogin implements UserBiz.
func (u *userBiz) UserLogin(ctx context.Context, req model.UserLoginRequest) (*model.UserInfo, error) {
	// 1. 根据用户名查询用户是否存在
	user, err := u.store.User().GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errno.ErrDatabase // 数据库查询错误
	}
	if user == nil {
		return nil, errno.ErrUserNotFound // 用户不存在
	}

	// 2. 校验密码（输入密码与数据库存储的加密密码比对）
	if err := authn.Compare(user.Password, req.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect // 密码错误
	}

	userInfo := &model.UserInfo{
		UserID:    user.UserID,
		Username:  user.Username,
		UpdatedAt: user.UpdatedAt,
	}
	// 4. 构造登录响应（可根据需求补充其他用户信息）
	// return &model.UserLoginResponse{
	// 	User:  *userInfo,
	// 	Token: tokenString,
	// }, nil

	return userInfo, nil

}

// CreteUser implements UserBiz.
func (u *userBiz) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserInfo, error) {

	// 判断用户名是否已存在
	existingUser, err := u.store.User().GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errno.ErrUserAlreadyExist
	}
	// 创建用户
	user := &model.User{
		UserID:    genid.GenerateUserID(),
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Phone:     req.Phone,
		NickName:  req.Nickname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// 对密码进行加密处理
	// 采用 bcrypt 加密密码
	password, err := authn.Encrypt(req.Password)
	if err != nil {
		return nil, errno.ErrEncryptPassword
	}
	user.Password = password
	// 保存用户到数据库
	if err := u.store.User().Create(ctx, user); err != nil {
		return nil, err
	}
	userInfo := &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Nickname: user.NickName,
	}
	return userInfo, nil

}

// DeleteUser implements UserBiz.
func (u *userBiz) DeleteUser(ctx context.Context, userID string) error {
	//  判断用户是否存在
	_, err := u.store.User().GetByID(ctx, userID)
	if err != nil {
		return err
	}
	err = u.store.User().Delete(ctx, userID)
	if err != nil {
		return err
	}
	return nil

}

// GetUserByID implements UserBiz.
func (u *userBiz) GetUserByID(ctx context.Context, userID string) (*model.UserInfo, error) {
	user, err := u.store.User().GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userInfo := &model.UserInfo{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Nickname: user.NickName,
	}
	return userInfo, nil
}

// GetUserByUsername implements UserBiz.
func (u *userBiz) GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error) {
	user, err := u.store.User().GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	count, _, err := u.store.Post().GetByUserID(ctx, user.UserID, 1, 10)
	if err != nil {
		return nil, err
	}
	//  TODO 查询blog信息,查询当前用户的blog信息
	userinfo := &model.UserInfo{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Nickname:  user.NickName,
		BlogTotal: count,
	}
	return userinfo, nil
}

// ListUsers implements UserBiz.
func (u *userBiz) ListUsers(ctx context.Context, page int, pageSize int) (*model.ListUserResponse, error) {
	userList, totalCount, err := u.store.User().List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	userInfoList := make([]model.UserInfo, 0, len(userList))
	// 查询每个用户的博客数量
	// 遍历userlist 查询每个用户的博客数量
	eg, ctx := errgroup.WithContext(ctx)
	// 设置最大并发数
	eg.SetLimit(10)
	var count int64
	for _, user := range userList {

		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			count, _, err = u.store.Post().GetByUserID(ctx, user.UserID, 1, 10)
			if err != nil {
				return err
			}
			log.Log.Error("查询用户博客失败", zap.Error(err))
			userInfoList = append(userInfoList, model.UserInfo{
				UserID:    user.UserID,
				Username:  user.Username,
				Email:     user.Email,
				Phone:     user.Phone,
				Nickname:  user.NickName,
				BlogTotal: count,
			})
			return nil
		})

	}
	if err := eg.Wait(); err != nil {
		log.Log.Error(err.Error())
		return nil, err
	}

	userListResponse := &model.ListUserResponse{
		TotalCount: totalCount, // TODO 从数据库中获取
		User:       userInfoList,
		HasMore:    totalCount > int64(page*pageSize), // 假设每页10条记录，判断是否还有更多记录,
	}

	return userListResponse, nil
}

// UpdateUser implements UserBiz.
func (u *userBiz) UpdateUser(ctx context.Context, updateUser *model.UpdateUser) error {
	// 判断当前用户是否存在

	userexist, err := u.store.User().GetByID(ctx, updateUser.UserID)
	if err != nil {
		return errno.ErrNotFound
	}
	user := &model.User{
		UserID:    userexist.UserID,
		UpdatedAt: time.Now(),
	}
	if updateUser.Email != "" && updateUser.Email != userexist.Email {
		user.Email = updateUser.Email
	}
	if updateUser.Phone != "" && updateUser.Phone != userexist.Phone {
		user.Phone = updateUser.Phone
	}
	if updateUser.Nickname != "" && updateUser.Nickname != userexist.NickName {
		user.NickName = updateUser.Nickname
	}

	err = u.store.User().Update(ctx, user)
	if err != nil {
		return err
	}

	return nil

}
