package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

type UserStore interface {
	// Create 创建用户
	Create(ctx context.Context, user *model.User) error
	// GetByID 根据 ID 获取用户
	GetByID(ctx context.Context, id uint) (*model.User, error)
	// GetByUUID 根据 UUID 获取用户
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	// Update 更新用户
	Update(ctx context.Context, user *model.User) error
	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
	// List 获取用户列表
	List(ctx context.Context, page, pageSize int) ([]*model.User, error)
}

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
	return u.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据 ID 获取用户
func (u *users) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUUID 根据 UUID 获取用户
func (u *users) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("user_id = ?", uuid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (u *users) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (u *users) Update(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (u *users) Delete(ctx context.Context, id uint) error {
	return u.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (u *users) List(ctx context.Context, page, pageSize int) ([]*model.User, error) {
	var users []*model.User
	offset := (page - 1) * pageSize
	err := u.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
