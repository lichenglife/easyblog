package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

// users 实现 UserStore 接口
type users struct {
	db *gorm.DB
}

// NewUsers 创建 users 实例
func NewUsers(db *gorm.DB) UserStore {
	return &users{db: db}
}

// Create 创建用户
func (u *users) Create(ctx context.Context, user *model.User) error {

	return nil
}

// GetByID 根据 ID 获取用户
func (u *users) GetByID(ctx context.Context, id uint) (*model.User, error) {

	return nil, nil
}

// GetByUsername 根据用户名获取用户
func (u *users) GetByUsername(ctx context.Context, username string) (*model.User, error) {

	return nil, nil
}

// Update 更新用户
func (u *users) Update(ctx context.Context, user *model.User) error {

	return nil
}

// Delete 删除用户
func (u *users) Delete(ctx context.Context, id uint) error {

	return nil
}

// List 获取用户列表
func (u *users) List(ctx context.Context, page, pageSize int) ([]*model.User, error) {

	return nil, nil
}
