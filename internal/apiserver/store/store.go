package store

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/gorm"
)

// Factory 存储层工厂接口
type IStore interface {
	// Post() PostStore
	User() UserStore

	Post() PostStore

	Close() error
}

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
	Delete(ctx context.Context, username string ) error
	// List 获取用户列表
	List(ctx context.Context, page, pageSize int) (users []*model.User, totalCount int64, err error)
}

type PostStore interface {
	// Create 创建帖子
	Create(ctx context.Context, post *model.Post) error
	// GetByID 根据 ID 获取帖子
	GetByID(ctx context.Context, postID string) (*model.Post, error)
	// Update 更新帖子
	Update(ctx context.Context, post *model.Post) error
	// Delete 删除帖子
	Delete(ctx context.Context, postID string) error
	// List 获取帖子列表
	List(ctx context.Context, page, pageSize int) ([]*model.Post, error)
	// GetByUserID 根据用户 ID 获取帖子列表
	GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.Post, error)
	// GetByPostID 根据帖子 ID 获取帖子
	GetByPostID(ctx context.Context, postID string) (*model.Post, error)
}

// dataStore 结构体, 实现Factory 接口
type dataStore struct {
	db *gorm.DB
}

// NewFactory 创建存储层工厂
func NewIStore(db *gorm.DB) IStore {
	return &dataStore{db: db}
}

// User() UserStore
func (ds *dataStore) User() UserStore {
	return NewUsers(ds.db)
}

// Post() PostStore

func (ds *dataStore) Post() PostStore {
	return NewPosts(ds.db)
}

func (ds *dataStore) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
