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

	return u.db.Create(&user).Error
}

// GetByID 根据 ID 获取用户
func (u *users) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("userID = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

// GetByUsername 根据用户名获取用户
func (u *users) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (u *users) Update(ctx context.Context, user *model.User) error {

	return u.db.Model(&user).Where("username = ?", user.Username).Updates(map[string]interface{}{
		"password": user.Password,
	}).Error

}

// Delete 删除用户
func (u *users) Delete(ctx context.Context, username string) error {
	if err := u.db.Delete("username = ?", username).Error; err != nil {
		return err
	}
	return nil
}

// List 获取用户列表
func (u *users) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	// 获取总条数
	var count int64
	if err := u.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var users []*model.User
	if err := u.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, count, nil
}
