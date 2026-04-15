package model

import "time"

// Tag 标签模型
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(30);not null;uniqueIndex:idx_tags_name;comment:标签名称" json:"name"`
	Slug      string    `gorm:"column:slug;type:varchar(30);not null;uniqueIndex:idx_tags_slug;comment:标签别名" json:"slug"`
	PostCount uint      `gorm:"column:post_count;type:int unsigned;not null;default:0;comment:文章数" json:"postCount"`
	Status    int       `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:状态：0-禁用 1-启用" json:"status"`
	CreateAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
	UpdateAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateAt"`
}

// TableName 表名
func (Tag) TableName() string { return "tags" }

// 创建标签请求结构
type CreateTagRequest struct {
	Name   string `json:"name" binding:"required,max=30"`
	Slug   string `json:"slug" binding:"required,max=30"`
}

// 更新标签请求结构
type UpdateTagRequest struct {
	ID     uint   `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required,max=30"`
	Slug   string `json:"slug" binding:"required,max=30"`
	Status int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// 标签响应结构
type TagInfo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	PostCount uint      `json:"postCount"`
	Status    int       `json:"status"`
	CreateAt  string    `json:"createAt"`
	UpdateAt  string    `json:"updateAt"`
}

// 标签列表响应
type TagListResponse struct {
	TotalCount int64     `json:"totalCount"`
	HasMore    bool      `json:"hasMore"`
	Tags       []TagInfo `json:"tags"`
}
