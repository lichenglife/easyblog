package commentv1

import (
	"context"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"gorm.io/gorm"
)

type CommentBiz interface {
	// CreateComment 创建评论
	CreateComment(ctx context.Context, userID uint, req *model.CreateCommentRequest) (*model.CommentInfo, error)
	// GetCommentByID 根据 ID 获取评论
	GetCommentByID(ctx context.Context, id uint) (*model.CommentInfo, error)
	// UpdateComment 更新评论
	UpdateComment(ctx context.Context, userID uint, req *model.UpdateCommentRequest) (*model.CommentInfo, error)
	// DeleteComment 删除评论
	DeleteComment(ctx context.Context, id uint) error
	// ListComments 获取文章评论列表（顶级评论）
	ListComments(ctx context.Context, postID uint, page, pageSize int) (*model.CommentListResponse, error)
	// ListReplies 获取楼中楼回复列表
	ListReplies(ctx context.Context, parentID uint, page, pageSize int) (*model.CommentListResponse, error)
	// LikeComment 点赞评论
	LikeComment(ctx context.Context, userID, commentID uint) error
	// UnlikeComment 取消点赞评论
	UnlikeComment(ctx context.Context, userID, commentID uint) error
}

type commentBiz struct {
	store store.CommentStore
	likeStore store.LikeStore
}

var _ CommentBiz = (*commentBiz)(nil)

func NewCommentBiz(store store.CommentStore, likeStore store.LikeStore) CommentBiz {
	return &commentBiz{
		store: store,
		likeStore: likeStore,
	}
}

// CreateComment 创建评论
func (b *commentBiz) CreateComment(ctx context.Context, userID uint, req *model.CreateCommentRequest) (*model.CommentInfo, error) {
	// 验证父评论是否存在（如果有）
	if req.ParentID > 0 {
		parent, err := b.store.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, errno.ErrCommentNotFound
		}
		if parent.PostID != req.PostID {
			return nil, errno.ErrInvalidParams.WithMessage("父评论不属于当前文章")
		}
		// 增加父评论的回复数
		if err := b.store.IncrementReplyCount(ctx, req.ParentID); err != nil {
			return nil, errno.ErrDatabase.WithMessage(err.Error())
		}
	}

	comment := &model.Comment{
		PostID:     req.PostID,
		UserID:     userID,
		ParentID:   req.ParentID,
		Content:    req.Content,
		IsEdited:   0,
		LikeCount:  0,
		ReplyCount: 0,
		Status:     1, // 默认正常
	}

	if err := b.store.Create(ctx, comment); err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.CommentInfo{
		ID:         comment.ID,
		PostID:     comment.PostID,
		UserID:     comment.UserID,
		ParentID:   comment.ParentID,
		Content:    comment.Content,
		IsEdited:   comment.IsEdited,
		LikeCount:  comment.LikeCount,
		ReplyCount: comment.ReplyCount,
		Status:     comment.Status,
		CreateAt:   comment.CreateAt.Format(time.RFC3339),
		UpdateAt:   comment.UpdateAt.Format(time.RFC3339),
	}, nil
}

// GetCommentByID 根据 ID 获取评论
func (b *commentBiz) GetCommentByID(ctx context.Context, id uint) (*model.CommentInfo, error) {
	comment, err := b.store.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrCommentNotFound
		}
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.CommentInfo{
		ID:         comment.ID,
		PostID:     comment.PostID,
		UserID:     comment.UserID,
		ParentID:   comment.ParentID,
		Content:    comment.Content,
		IsEdited:   comment.IsEdited,
		LikeCount:  comment.LikeCount,
		ReplyCount: comment.ReplyCount,
		Status:     comment.Status,
		CreateAt:   comment.CreateAt.Format(time.RFC3339),
		UpdateAt:   comment.UpdateAt.Format(time.RFC3339),
	}, nil
}

// UpdateComment 更新评论
func (b *commentBiz) UpdateComment(ctx context.Context, userID uint, req *model.UpdateCommentRequest) (*model.CommentInfo, error) {
	comment, err := b.store.GetByID(ctx, req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrCommentNotFound
		}
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	// 只能修改自己的评论
	if comment.UserID != userID {
		return nil, errno.ErrPermissionDenied
	}

	comment.Content = req.Content
	comment.IsEdited = 1

	if err := b.store.Update(ctx, comment); err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	return &model.CommentInfo{
		ID:         comment.ID,
		PostID:     comment.PostID,
		UserID:     comment.UserID,
		ParentID:   comment.ParentID,
		Content:    comment.Content,
		IsEdited:   comment.IsEdited,
		LikeCount:  comment.LikeCount,
		ReplyCount: comment.ReplyCount,
		Status:     comment.Status,
		CreateAt:   comment.CreateAt.Format(time.RFC3339),
		UpdateAt:   comment.UpdateAt.Format(time.RFC3339),
	}, nil
}

// DeleteComment 删除评论
func (b *commentBiz) DeleteComment(ctx context.Context, id uint) error {
	comment, err := b.store.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.ErrCommentNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 软删除：将状态设为 2-已删除
	comment.Status = 2
	if err := b.store.Update(ctx, comment); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 如果是回复评论，减少父评论的回复数
	if comment.ParentID > 0 {
		_ = b.store.DecrementReplyCount(ctx, comment.ParentID)
	}

	return nil
}

// ListComments 获取文章评论列表（顶级评论）
func (b *commentBiz) ListComments(ctx context.Context, postID uint, page, pageSize int) (*model.CommentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	comments, err := b.store.ListByPostID(ctx, postID, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	commentInfos := make([]model.CommentInfo, 0, len(comments))
	for _, c := range comments {
		commentInfos = append(commentInfos, model.CommentInfo{
			ID:         c.ID,
			PostID:     c.PostID,
			UserID:     c.UserID,
			ParentID:   c.ParentID,
			Content:    c.Content,
			IsEdited:   c.IsEdited,
			LikeCount:  c.LikeCount,
			ReplyCount: c.ReplyCount,
			Status:     c.Status,
			CreateAt:   c.CreateAt.Format(time.RFC3339),
			UpdateAt:   c.UpdateAt.Format(time.RFC3339),
		})
	}

	// 简单计算总数（实际项目中应该用 count 查询）
	totalCount := int64(len(commentInfos))
	hasMore := len(commentInfos) == pageSize

	return &model.CommentListResponse{
		TotalCount: totalCount,
		HasMore:    hasMore,
		Comments:   commentInfos,
	}, nil
}

// ListReplies 获取楼中楼回复列表
func (b *commentBiz) ListReplies(ctx context.Context, parentID uint, page, pageSize int) (*model.CommentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	comments, err := b.store.ListByParentID(ctx, parentID, page, pageSize)
	if err != nil {
		return nil, errno.ErrDatabase.WithMessage(err.Error())
	}

	commentInfos := make([]model.CommentInfo, 0, len(comments))
	for _, c := range comments {
		commentInfos = append(commentInfos, model.CommentInfo{
			ID:         c.ID,
			PostID:     c.PostID,
			UserID:     c.UserID,
			ParentID:   c.ParentID,
			Content:    c.Content,
			IsEdited:   c.IsEdited,
			LikeCount:  c.LikeCount,
			ReplyCount: c.ReplyCount,
			Status:     c.Status,
			CreateAt:   c.CreateAt.Format(time.RFC3339),
			UpdateAt:   c.UpdateAt.Format(time.RFC3339),
		})
	}

	totalCount := int64(len(commentInfos))
	hasMore := len(commentInfos) == pageSize

	return &model.CommentListResponse{
		TotalCount: totalCount,
		HasMore:    hasMore,
		Comments:   commentInfos,
	}, nil
}

// LikeComment 点赞评论
func (b *commentBiz) LikeComment(ctx context.Context, userID, commentID uint) error {
	// 检查评论是否存在
	_, err := b.store.GetByID(ctx, commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.ErrCommentNotFound
		}
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 检查是否已点赞
	liked, _ := b.likeStore.IsLiked(ctx, userID, commentID, 2) // 2=评论
	if liked {
		return errno.ErrAlreadyLiked
	}

	// 创建点赞记录
	like := &model.Like{
		UserID:     userID,
		TargetID:   commentID,
		TargetType: 2, // 2=评论
	}
	if err := b.likeStore.Create(ctx, like); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 增加评论点赞数
	if err := b.store.IncrementLikeCount(ctx, commentID); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}

// UnlikeComment 取消点赞评论
func (b *commentBiz) UnlikeComment(ctx context.Context, userID, commentID uint) error {
	// 删除点赞记录
	if err := b.likeStore.Delete(ctx, userID, commentID, 2); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	// 减少评论点赞数
	if err := b.store.DecrementLikeCount(ctx, commentID); err != nil {
		return errno.ErrDatabase.WithMessage(err.Error())
	}

	return nil
}
