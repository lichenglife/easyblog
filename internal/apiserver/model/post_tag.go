package model

import "time"

// PostTag 文章 - 标签关联模型
type PostTag struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	PostID   uint      `gorm:"column:post_id;type:bigint unsigned;not null;uniqueIndex:idx_post_tags_post_id;comment:文章 ID" json:"postId"`
	TagID    uint      `gorm:"column:tag_id;type:bigint unsigned;not null;index:idx_post_tags_tag_id;comment:标签 ID" json:"tagId"`
	CreateAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
}

// TableName 表名
func (PostTag) TableName() string { return "post_tags" }
