package biz

import (
	commentv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/comment"
	likev1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/like"
	postv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/post"
	userv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/user"
	categoryv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/category"
	tagv1 "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/tag"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/auth"
	"go.uber.org/zap"
)

type IBiz interface {
	// 用户业务接口 V1 版本
	UserV1() userv1.UserBiz
	// 博客业务接口 V1 版本
	PostV1() postv1.PostBiz
	// 评论业务接口 V1 版本
	CommentV1() commentv1.CommentBiz
	// 点赞业务接口 V1 版本
	LikeV1() likev1.LikeBiz
	// 分类业务接口 V1 版本
	CategoryV1() categoryv1.CategoryBiz
	// 标签业务接口 V1 版本
	TagV1() tagv1.TagBiz
}

// biz 实现 IBiz 接口
type biz struct {
	// 存储层的业务逻辑
	store  store.IStore
	jwt    *auth.JWT
	logger *zap.Logger
}

// PostV1 implements IBiz.
func (b *biz) PostV1() postv1.PostBiz {
	return postv1.NewPostBiz(b.store.Post(), b.store.PostTag())
}

// UserV1 implements IBiz.
func (b *biz) UserV1() userv1.UserBiz {
	return userv1.NewUserBiz(b.store.User(), b.jwt)
}

// CommentV1 implements IBiz.
func (b *biz) CommentV1() commentv1.CommentBiz {
	return commentv1.NewCommentBiz(b.store.Comment(), b.store.Like())
}

// LikeV1 implements IBiz.
func (b *biz) LikeV1() likev1.LikeBiz {
	return likev1.NewLikeBiz(b.store.Like(), b.store.Post())
}

// CategoryV1 implements IBiz.
func (b *biz) CategoryV1() categoryv1.CategoryBiz {
	return categoryv1.NewCategoryBiz(b.logger, b.store)
}

// TagV1 implements IBiz.
func (b *biz) TagV1() tagv1.TagBiz {
	return tagv1.NewTagBiz(b.logger, b.store)
}

// NewBiz 创建业务逻辑层实例
func NewBiz(store store.IStore, jwt *auth.JWT, logger *zap.Logger) IBiz {
	return &biz{store: store, jwt: jwt, logger: logger}
}
