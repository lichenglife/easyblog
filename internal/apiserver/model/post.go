package model

import "time"

// Post 博客模型
type Post struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"column:user_id;type:bigint unsigned;not null;index:idx_posts_user_id;comment:作者 ID" json:"userId"`
	PostID        string    `gorm:"column:post_id;type:varchar(36);not null;uniqueIndex:idx_posts_post_id;comment:文章唯一 ID" json:"postId"`
	Title         string    `gorm:"column:title;type:varchar(255);not null;comment:标题" json:"title"`
	Content       string    `gorm:"column:content;type:longtext;not null;comment:内容（Markdown）" json:"content"`
	Summary       string    `gorm:"column:summary;type:varchar(500);comment:摘要" json:"summary"`
	CoverImage    string    `gorm:"column:cover_image;type:varchar(255);comment:封面图 URL" json:"coverImage"`
	CategoryID    uint      `gorm:"column:category_id;type:bigint unsigned;comment:分类 ID" json:"categoryId"`
	Status        int       `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:状态：0-草稿 1-已发布 2-已下架" json:"status"`
	IsTop         int       `gorm:"column:is_top;type:tinyint unsigned;not null;default:0;comment:是否置顶：0-否 1-是" json:"isTop"`
	ViewCount     uint      `gorm:"column:view_count;type:int unsigned;not null;default:0;comment:阅读量" json:"viewCount"`
	LikeCount     uint      `gorm:"column:like_count;type:int unsigned;not null;default:0;comment:点赞数" json:"likeCount"`
	CommentCount  uint      `gorm:"column:comment_count;type:int unsigned;not null;default:0;comment:评论数" json:"commentCount"`
	CreateAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
	UpdateAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateAt"`
}

// TableName 表名
func (Post) TableName() string { return "posts" }

// 创建文章请求结构
type CreatePostRequest struct {
	Title      string `json:"title" binding:"required,max=255"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary" binding:"omitempty,max=500"`
	CoverImage string `json:"coverImage" binding:"omitempty"`
	CategoryID uint   `json:"categoryId" binding:"omitempty"`
	Status     int    `json:"status" binding:"omitempty,oneof=0 1 2"` // 0-草稿 1-已发布 2-已下架
	IsTop      int    `json:"isTop" binding:"omitempty,oneof=0 1"`
	TagIDs     []uint `json:"tagIds" binding:"omitempty"` // 标签 ID 列表
}

// 更新文章请求结构
type UpdatePostRequest struct {
	ID         uint   `json:"id" binding:"required"`
	Title      string `json:"title" binding:"required,max=255"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary" binding:"omitempty,max=500"`
	CoverImage string `json:"coverImage" binding:"omitempty"`
	CategoryID uint   `json:"categoryId" binding:"omitempty"`
	Status     int    `json:"status" binding:"omitempty,oneof=0 1 2"`
	IsTop      int    `json:"isTop" binding:"omitempty,oneof=0 1"`
	TagIDs     []uint `json:"tagIds" binding:"omitempty"` // 标签 ID 列表
}

// 文章响应结构
type PostInfo struct {
	ID           uint   `json:"id"`
	PostID       string `json:"postId"`
	UserID       uint   `json:"userId"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Summary      string `json:"summary"`
	CoverImage   string `json:"coverImage"`
	CategoryID   uint   `json:"categoryId"`
	Status       int    `json:"status"`
	IsTop        int    `json:"isTop"`
	ViewCount    uint   `json:"viewCount"`
	LikeCount    uint   `json:"likeCount"`
	CommentCount uint   `json:"commentCount"`
	CreateAt     string `json:"createAt"`
	UpdateAt     string `json:"updateAt"`
	Tags         []uint `json:"tags,omitempty"` // 标签 ID 列表
}

// 查询列表请求
type ListPostResponse struct {
	TotalCount int64      `json:"totalCount"`
	HasMore    bool       `json:"hasMore"`
	Posts      []PostInfo `json:"posts"`
}

// 按分类筛选请求
type ListPostsByCategoryRequest struct {
	CategoryID uint `json:"categoryId" form:"categoryId"`
	Page       int  `json:"page" form:"page"`
	PageSize   int  `json:"pageSize" form:"pageSize"`
}

// 按标签筛选请求
type ListPostsByTagRequest struct {
	TagID    uint `json:"tagId" form:"tagId"`
	Page     int  `json:"page" form:"page"`
	PageSize int  `json:"pageSize" form:"pageSize"`
}

// 搜索文章请求
type SearchPostsRequest struct {
	Keyword  string `json:"keyword" form:"keyword"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}
