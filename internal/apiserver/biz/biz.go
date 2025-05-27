package biz

import (
	"context"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// UserBiz 用户业务接口
type UserBiz interface {
	// Create 创建用户
	CreteUser(ctx context.Context, req *model.User) (*model.User, error)
	// GetByID 根据 ID 获取用户
	GetUserByID(ctx context.Context, userID string) (*model.UserInfo, error)
	// GetByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error)
	// Update 更新用户
	UpdateUser(ctx context.Context, username string, user *model.UpdateUser) error
	// Delete 删除用户
	DeleteUser(ctx context.Context, username string) error
	// List 获取用户列表
	ListUsers(ctx context.Context, page, pageSize int) (*model.ListUserResponse, error)

	ResetUserPassword(ctx context.Context, username, oldPassword, newPassword string) error

	Login(ctx context.Context, username, password string) (*model.UserInfo, error)
}

type PostBiz interface {
	// Create 创建帖子
	CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error)
	// GetByID 根据 ID 获取帖子
	//GetPostByID(ctx context.Context, id uint) (*model.Post, error)
	// Update 更新帖子
	UpdatePost(ctx context.Context, post *model.UpdatePostRequest) error
	// Delete 删除帖子
	DeletePostByPostID(ctx context.Context,  postID string) error
	// List 获取帖子列表
	ListPosts(ctx context.Context, page, pageSize int) (*model.ListPostResponse, error)
	// GetByUserID 根据用户 ID 获取帖子列表
	GetPostsByUserID(ctx context.Context, userID string, page, pageSize int) (int,*model.ListPostResponse, error)
	// GetByPostID 根据帖子 ID 获取帖子
	GetPostByPostID(ctx context.Context, postID string) (*model.Post, error)
}

// NewPostBiz 实例化postBiz对象
func NewPostBiz(logger *log.Logger, store store.IStore) PostBiz {

	return &postBiz{
		logger: logger,
		store:  store,
	}
}

func NewUserBiz(logger *log.Logger, store store.IStore) UserBiz {
	return &userBiz{
		logger: logger,
		store:  store,
	}
}
