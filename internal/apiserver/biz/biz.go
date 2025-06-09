package biz

import (
	postv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/post"
	userv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/user"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
)

type IBiz interface {
	// 用户业务接口V1版本
	UserV1() userv1.UserBiz
	// 博客业务接口V1版本
	PostV1() postv1.PostBiz
}

// biz 实现 IBiz 接口
type biz struct {
	// 存储层的业务逻辑
	store store.IStore
}

// PostV1 implements IBiz.
func (b *biz) PostV1() postv1.PostBiz {
	return   postv1.NewPostBiz(b.store.Post())
}

// UserV1 implements IBiz.
func (b *biz) UserV1() userv1.UserBiz {
	return userv1.NewUserBiz(b.store.User())
}

// NewBiz 创建业务逻辑层实例
func NewBiz(store store.IStore) IBiz {
	return &biz{store: store}
}
