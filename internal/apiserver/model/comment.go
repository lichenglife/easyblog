package model

import "time"

// Comment 评论模型
type Comment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PostID      uint      `gorm:"column:post_id;type:bigint unsigned;not null;index:idx_comments_post_id;comment:文章 ID" json:"postId"`
	UserID      uint      `gorm:"column:user_id;type:bigint unsigned;not null;index:idx_comments_user_id;comment:评论者 ID" json:"userId"`
	ParentID    uint      `gorm:"column:parent_id;type:bigint unsigned;default:0;index:idx_comments_parent_id;comment:父评论 ID（楼中楼）" json:"parentId"`
	Content     string    `gorm:"column:content;type:text;not null;comment:评论内容" json:"content"`
	IsEdited    int       `gorm:"column:is_edited;type:tinyint unsigned;not null;default:0;comment:是否已编辑：0-否 1-是" json:"isEdited"`
	LikeCount   uint      `gorm:"column:like_count;type:int unsigned;not null;default:0;comment:点赞数" json:"likeCount"`
	ReplyCount  uint      `gorm:"column:reply_count;type:int unsigned;not null;default:0;comment:回复数" json:"replyCount"`
	Status      int       `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:状态：0-待审核 1-正常 2-已删除" json:"status"`
	CreateAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
	UpdateAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateAt"`
}

// TableName 表名
func (Comment) TableName() string { return "comments" }

// 创建评论请求结构
type CreateCommentRequest struct {
	PostID   uint   `json:"postId" binding:"required"`
	ParentID uint   `json:"parentId" binding:"omitempty"`
	Content  string `json:"content" binding:"required,max=1000"`
}

// 更新评论请求结构
type UpdateCommentRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Content string `json:"content" binding:"required,max=1000"`
}

// 评论响应结构
type CommentInfo struct {
	ID         uint      `json:"id"`
	PostID     uint      `json:"postId"`
	UserID     uint      `json:"userId"`
	ParentID   uint      `json:"parentId"`
	Content    string    `json:"content"`
	IsEdited   int       `json:"isEdited"`
	LikeCount  uint      `json:"likeCount"`
	ReplyCount uint      `json:"replyCount"`
	Status     int       `json:"status"`
	CreateAt   string    `json:"createAt"`
	UpdateAt   string    `json:"updateAt"`
}

// 评论列表响应
type CommentListResponse struct {
	TotalCount int64         `json:"totalCount"`
	HasMore    bool          `json:"hasMore"`
	Comments   []CommentInfo `json:"comments"`
}
