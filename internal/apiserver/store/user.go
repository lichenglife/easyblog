package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserStore interface {
	// Create 创建用户
	Create(ctx context.Context, user *model.User) error
	// GetByID 根据 ID 获取用户
	GetByID(ctx context.Context, userID string) (*model.User, error)
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	// Update 更新用户
	Update(ctx context.Context, user *model.User) error
	// Delete 删除用户
	Delete(ctx context.Context, username string) error
	// List 获取用户列表
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
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

	return u.db.Create(user).Error
}

// GetByID 根据 ID 获取用户
func (u *users) GetByID(ctx context.Context, userID string) (*model.User, error) {
	var user *model.User
	err := u.db.WithContext(ctx).Where("userID =?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (u *users) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // 用户不存在返回 nil
	}
	return &user, err
}

// Update 更新用户
func (u *users) Update(ctx context.Context, user *model.User) error {
	// 更新 user对象
	if err := u.db.WithContext(ctx).Updates(user).Error; err != nil {
		log.Log.Error("failed to update object", zap.Error(err))
		return err
	}
	return nil

}

// Delete 删除用户
func (u *users) Delete(ctx context.Context, username string) error {
	return u.db.WithContext(ctx).Delete(&model.User{}).Where("username = ?", username).Error
}

// List 获取用户列表
func (u *users) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64
	err := u.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = u.db.WithContext(ctx).Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, 0, nil
}
