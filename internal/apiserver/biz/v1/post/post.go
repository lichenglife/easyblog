package biz

import (
	"context"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/util"
)

type PostBiz interface {
	// Create 创建文章
	CreatePost(ctx context.Context, userID uint, req *model.CreatePostRequest) (*model.PostInfo, error)
	// GetByID 根据 ID 获取文章
	GetPostByID(ctx context.Context, id uint) (*model.PostInfo, error)
	// Update 更新文章
	UpdatePost(ctx context.Context, id uint, req *model.UpdatePostRequest) error
	// Delete 删除文章
	DeletePost(ctx context.Context, id uint) error
	// List 获取文章列表
	ListPosts(ctx context.Context, page, pageSize int) (*model.ListPostResponse, error)
	// GetByUserID 根据用户 ID 获取文章列表
	GetPostsByUserID(ctx context.Context, userID uint, page, pageSize int) (*model.ListPostResponse, error)
	// GetByPostID 根据文章 ID 获取文章
	GetPostByPostID(ctx context.Context, postID string) (*model.PostInfo, error)
	// SearchPosts 搜索文章
	SearchPosts(ctx context.Context, keyword string, page, pageSize int) (*model.ListPostResponse, error)
	// GetMyPosts 获取我的文章列表
	GetMyPosts(ctx context.Context, userID uint, page, pageSize int) (*model.ListPostResponse, error)
	// ListByCategoryID 根据分类 ID 获取文章列表
	ListByCategoryID(ctx context.Context, categoryID uint, page, pageSize int) (*model.ListPostResponse, error)
	// ListByTagID 根据标签 ID 获取文章列表
	ListByTagID(ctx context.Context, tagID uint, page, pageSize int) (*model.ListPostResponse, error)
	// TopPost 置顶文章
	TopPost(ctx context.Context, id uint) error
	// UnTopPost 取消置顶文章
	UnTopPost(ctx context.Context, id uint) error
}

// NewPostBiz 实例化 postBiz 对象
func NewPostBiz(postStore store.PostStore, postTagStore store.PostTagStore) PostBiz {
	return &postBiz{
		postStore:    postStore,
		postTagStore: postTagStore,
	}
}

// postBiz 实现了 post 业务层接口
type postBiz struct {
	postStore    store.PostStore
	postTagStore store.PostTagStore
}

// CreatePost implements PostBiz.
func (p *postBiz) CreatePost(ctx context.Context, userID uint, req *model.CreatePostRequest) (*model.PostInfo, error) {
	post := &model.Post{
		PostID:       util.GenerateUUID(),
		UserID:       userID,
		Title:        req.Title,
		Content:      req.Content,
		Summary:      req.Summary,
		CoverImage:   req.CoverImage,
		CategoryID:   req.CategoryID,
		Status:       req.Status,
		IsTop:        req.IsTop,
		ViewCount:    0,
		LikeCount:    0,
		CommentCount: 0,
	}

	// 使用事务创建文章和标签关联
	if err := p.postStore.Create(ctx, post); err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	// 处理标签关联
	if len(req.TagIDs) > 0 {
		postTags := make([]*model.PostTag, 0, len(req.TagIDs))
		for _, tagID := range req.TagIDs {
			postTags = append(postTags, &model.PostTag{
				PostID: post.ID,
				TagID:  tagID,
			})
		}
		if err := p.postTagStore.CreateBatch(ctx, postTags); err != nil {
			return nil, errno.ErrDatabase.WithMessage(err.Error())
		}
	}

	return &model.PostInfo{
		ID:           post.ID,
		PostID:       post.PostID,
		UserID:       post.UserID,
		Title:        post.Title,
		Content:      post.Content,
		Summary:      post.Summary,
		CoverImage:   post.CoverImage,
		CategoryID:   post.CategoryID,
		Status:       post.Status,
		IsTop:        post.IsTop,
		ViewCount:    post.ViewCount,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		CreateAt:     post.CreateAt.Format(time.RFC3339),
		UpdateAt:     post.UpdateAt.Format(time.RFC3339),
		Tags:         req.TagIDs,
	}, nil
}

// DeletePost implements PostBiz.
func (p *postBiz) DeletePost(ctx context.Context, id uint) error {
	// 检查文章是否存在
	_, err := p.postStore.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrPostNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 先删除标签关联
	if err := p.postTagStore.DeleteByPostID(ctx, id); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 再删除文章
	if err := p.postStore.Delete(ctx, id); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}
	return nil
}

// GetPostByID implements PostBiz.
func (p *postBiz) GetPostByID(ctx context.Context, id uint) (*model.PostInfo, error) {
	post, err := p.postStore.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errno.ErrPostNotFound
		}
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	// 获取标签 ID 列表
	tagIDs, err := p.postTagStore.ListTagIDsByPostID(ctx, id)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.PostInfo{
		ID:           post.ID,
		PostID:       post.PostID,
		UserID:       post.UserID,
		Title:        post.Title,
		Content:      post.Content,
		Summary:      post.Summary,
		CoverImage:   post.CoverImage,
		CategoryID:   post.CategoryID,
		Status:       post.Status,
		IsTop:        post.IsTop,
		ViewCount:    post.ViewCount,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		CreateAt:     post.CreateAt.Format(time.RFC3339),
		UpdateAt:     post.UpdateAt.Format(time.RFC3339),
		Tags:         tagIDs,
	}, nil
}

// GetPostByPostID implements PostBiz.
func (p *postBiz) GetPostByPostID(ctx context.Context, postID string) (*model.PostInfo, error) {
	post, err := p.postStore.GetByPostID(ctx, postID)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errno.ErrPostNotFound
		}
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.PostInfo{
		ID:           post.ID,
		PostID:       post.PostID,
		UserID:       post.UserID,
		Title:        post.Title,
		Content:      post.Content,
		Summary:      post.Summary,
		CoverImage:   post.CoverImage,
		CategoryID:   post.CategoryID,
		Status:       post.Status,
		IsTop:        post.IsTop,
		ViewCount:    post.ViewCount,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		CreateAt:     post.CreateAt.Format(time.RFC3339),
		UpdateAt:     post.UpdateAt.Format(time.RFC3339),
	}, nil
}

// GetPostsByUserID implements PostBiz.
func (p *postBiz) GetPostsByUserID(ctx context.Context, userID uint, page int, pageSize int) (*model.ListPostResponse, error) {
	posts, err := p.postStore.ListByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	postInfos, err := p.buildPostInfos(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &model.ListPostResponse{
		TotalCount: int64(len(postInfos)),
		HasMore:    len(postInfos) == pageSize,
		Posts:      postInfos,
	}, nil
}

// ListPosts implements PostBiz.
func (p *postBiz) ListPosts(ctx context.Context, page int, pageSize int) (*model.ListPostResponse, error) {
	posts, err := p.postStore.List(ctx, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	postInfos, err := p.buildPostInfos(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &model.ListPostResponse{
		TotalCount: int64(len(postInfos)),
		HasMore:    len(postInfos) == pageSize,
		Posts:      postInfos,
	}, nil
}

// UpdatePost implements PostBiz.
func (p *postBiz) UpdatePost(ctx context.Context, id uint, req *model.UpdatePostRequest) error {
	post, err := p.postStore.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrPostNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	post.Title = req.Title
	post.Content = req.Content
	post.Summary = req.Summary
	post.CoverImage = req.CoverImage
	post.CategoryID = req.CategoryID
	post.Status = req.Status
	post.IsTop = req.IsTop

	if err := p.postStore.Update(ctx, post); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 处理标签关联：先删除旧的，再创建新的
	if err := p.postTagStore.DeleteByPostID(ctx, id); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	if len(req.TagIDs) > 0 {
		postTags := make([]*model.PostTag, 0, len(req.TagIDs))
		for _, tagID := range req.TagIDs {
			postTags = append(postTags, &model.PostTag{
				PostID: post.ID,
				TagID:  tagID,
			})
		}
		if err := p.postTagStore.CreateBatch(ctx, postTags); err != nil {
			return errno.ErrDatabase.WithMessage(err.Error())
		}
	}

	return nil
}

var _ PostBiz = (*postBiz)(nil)

// SearchPosts 搜索文章
func (p *postBiz) SearchPosts(ctx context.Context, keyword string, page, pageSize int) (*model.ListPostResponse, error) {
	posts, err := p.postStore.SearchPosts(ctx, keyword, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	postInfos, err := p.buildPostInfos(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &model.ListPostResponse{
		TotalCount: int64(len(postInfos)),
		HasMore:    len(postInfos) == pageSize,
		Posts:      postInfos,
	}, nil
}

// GetMyPosts 获取我的文章列表
func (p *postBiz) GetMyPosts(ctx context.Context, userID uint, page, pageSize int) (*model.ListPostResponse, error) {
	posts, err := p.postStore.ListByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	postInfos, err := p.buildPostInfos(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &model.ListPostResponse{
		TotalCount: int64(len(postInfos)),
		HasMore:    len(postInfos) == pageSize,
		Posts:      postInfos,
	}, nil
}

// ListByCategoryID 根据分类 ID 获取文章列表
func (p *postBiz) ListByCategoryID(ctx context.Context, categoryID uint, page, pageSize int) (*model.ListPostResponse, error) {
	posts, err := p.postStore.ListByCategoryID(ctx, categoryID, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	postInfos, err := p.buildPostInfos(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &model.ListPostResponse{
		TotalCount: int64(len(postInfos)),
		HasMore:    len(postInfos) == pageSize,
		Posts:      postInfos,
	}, nil
}

// ListByTagID 根据标签 ID 获取文章列表
func (p *postBiz) ListByTagID(ctx context.Context, tagID uint, page, pageSize int) (*model.ListPostResponse, error) {
	posts, err := p.postStore.ListByTagID(ctx, tagID, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	postInfos, err := p.buildPostInfos(ctx, posts)
	if err != nil {
		return nil, err
	}

	return &model.ListPostResponse{
		TotalCount: int64(len(postInfos)),
		HasMore:    len(postInfos) == pageSize,
		Posts:      postInfos,
	}, nil
}

// TopPost 置顶文章
func (p *postBiz) TopPost(ctx context.Context, id uint) error {
	_, err := p.postStore.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrPostNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	if err := p.postStore.SetTop(ctx, id); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// UnTopPost 取消置顶文章
func (p *postBiz) UnTopPost(ctx context.Context, id uint) error {
	_, err := p.postStore.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			return errno.ErrPostNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	if err := p.postStore.CancelTop(ctx, id); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// buildPostInfos 构建文章信息列表（包含标签）
func (p *postBiz) buildPostInfos(ctx context.Context, posts []*model.Post) ([]model.PostInfo, error) {
	var postInfos []model.PostInfo
	for _, post := range posts {
		// 获取标签 ID 列表
		tagIDs, err := p.postTagStore.ListTagIDsByPostID(ctx, post.ID)
		if err != nil {
			return nil, errno.ErrDatabase.WithMessage(err.Error())
		}

		postInfos = append(postInfos, model.PostInfo{
			ID:           post.ID,
			PostID:       post.PostID,
			UserID:       post.UserID,
			Title:        post.Title,
			Content:      post.Content,
			Summary:      post.Summary,
			CoverImage:   post.CoverImage,
			CategoryID:   post.CategoryID,
			Status:       post.Status,
			IsTop:        post.IsTop,
			ViewCount:    post.ViewCount,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			CreateAt:     post.CreateAt.Format(time.RFC3339),
			UpdateAt:     post.UpdateAt.Format(time.RFC3339),
			Tags:         tagIDs,
		})
	}
	return postInfos, nil
}
