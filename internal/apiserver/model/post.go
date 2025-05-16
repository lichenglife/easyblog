package model

import "time"

// Post 博客模型
type Post struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	UserID   string    `gorm:"column:userID;type:varchar(36);not null;uniqueIndex:user.userID;comment:用户唯一 ID" json:"userID"`
	PostID   string    `gorm:"column:postID;type:varchar(36);not null;uniqueIndex:post.postID;comment:帖子唯一 ID" json:"postID"`
	Content  string    `gorm:"column:content;type:varchar(255);not null;comment:内容" json:"content"`
	Title    string    `gorm:"column:title;type:varchar(255);not null;comment:标题" json:"title"`
	CreateAt time.Time `gorm:"column:createAt;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
	UpdateAt time.Time `gorm:"column:updateAt;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateAt"`
}

// TableName 表名
func (Post) TableName() string { return "post" }

// 创建帖子请求结构
type CreatePostRequest struct {
	UserID  string `json:"userID" binding:"required"`
	Content string `json:"content" binding:"required"`
	Title   string `json:"title" binding:"required"`
}

// 修改帖子请求结构
type UpdatePostRequest struct {
	ID      uint   `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
	Title   string `json:"title" binding:"required"`
}

// 查询帖子请求结构
type GetPostRequest struct {
	ID uint `json:"id" binding:"required"`
}

// 查询列表请求
type PageListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=5,max=10"`
}

type ListPostResponse struct {
	TotalCount int64  `json:"totalCount"`
	HasMore    bool   `json:"hasMore"`
	Posts      []Post `json:"posts"`
}
