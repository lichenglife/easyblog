package biz

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
)

// UserBiz 用户业务接口
type UserBiz interface {
	// Create 创建用户
	CreteUser(ctx context.Context, req *model.CreatePostRequest) (*model.UserInfo, error)
	// GetByID 根据 ID 获取用户
	GetUserByID(ctx context.Context, id uint) (*model.UserInfo, error)
	// GetByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error)
	// Update 更新用户
	UpdateUser(ctx context.Context, user *model.UpdateUser) error
	// Delete 删除用户
	DeleteUser(ctx context.Context, id uint) error
	// List 获取用户列表
	ListUsers(ctx context.Context, page, pageSize int) (*model.ListUserResponse, error)
}

func NewUserBiz(store store.UserStore) UserBiz {
	return &userBiz{
		store: store,
	}
}

var _ UserBiz = (*userBiz)(nil)

// userBiz 定义了用户业务逻辑层
type userBiz struct {
	store store.UserStore
}

// CreteUser implements UserBiz.
func (u *userBiz) CreteUser(ctx context.Context, req *model.CreatePostRequest) (*model.UserInfo, error) {
	panic("unimplemented")
}

// DeleteUser implements UserBiz.
func (u *userBiz) DeleteUser(ctx context.Context, id uint) error {
	panic("unimplemented")
}

// GetUserByID implements UserBiz.
func (u *userBiz) GetUserByID(ctx context.Context, id uint) (*model.UserInfo, error) {
	panic("unimplemented")
}

// GetUserByUsername implements UserBiz.
func (u *userBiz) GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error) {
	panic("unimplemented")
}

// ListUsers implements UserBiz.
func (u *userBiz) ListUsers(ctx context.Context, page int, pageSize int) (*model.ListUserResponse, error) {
	panic("unimplemented")
}

// UpdateUser implements UserBiz.
func (u *userBiz) UpdateUser(ctx context.Context, user *model.UpdateUser) error {
	panic("unimplemented")
}
