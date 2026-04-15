package model

import "time"

// Like 点赞模型
type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id;type:bigint unsigned;not null;uniqueIndex:idx_like_unique;comment:用户 ID" json:"userId"`
	TargetID  uint      `gorm:"column:target_id;type:bigint unsigned;not null;uniqueIndex:idx_like_unique;comment:目标 ID（文章或评论）" json:"targetId"`
	TargetType int      `gorm:"column:target_type;type:tinyint unsigned;not null;comment:目标类型：1-文章 2-评论" json:"targetType"`
	CreateAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
}

// TableName 表名
func (Like) TableName() string { return "likes" }

// 点赞请求结构
type CreateLikeRequest struct {
	TargetID   uint `json:"targetId" binding:"required"`
	TargetType int  `json:"targetType" binding:"required,oneof=1 2"` // 1-文章 2-评论
}

// 点赞响应结构
type LikeInfo struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	TargetID   uint      `json:"targetId"`
	TargetType int       `json:"targetType"`
	CreateAt   string    `json:"createAt"`
}

// 点赞状态响应
type LikeStatusResponse struct {
	IsLiked bool `json:"isLiked"`
}
