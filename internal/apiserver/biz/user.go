package biz

import (
	"context"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

var _ UserBiz = (*userBiz)(nil)

// userBiz 定义了用户业务逻辑层
type userBiz struct {
	logger *log.Logger
	store  store.UserStore
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
