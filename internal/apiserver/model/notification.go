package model

import "time"

// Notification 通知模型
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id;type:bigint unsigned;not null;index:idx_notifications_user_id;comment:接收通知的用户 ID" json:"userId"`
	Type      int       `gorm:"column:type;type:tinyint unsigned;not null;comment:通知类型：1-评论文章 2-回复评论 3-点赞文章 4-点赞评论" json:"type"`
	TargetID  uint      `gorm:"column:target_id;type:bigint unsigned;not null;index:idx_notifications_target_id;comment:目标对象 ID(评论 ID/文章 ID)" json:"targetId"`
	TargetType int      `gorm:"column:target_type;type:tinyint unsigned;not null;comment:目标对象类型：1-文章 2-评论" json:"targetType"`
	Content   string    `gorm:"column:content;type:varchar(500);comment:通知内容摘要" json:"content"`
	IsRead    int       `gorm:"column:is_read;type:tinyint unsigned;not null;default:0;comment:是否已读：0-未读 1-已读" json:"isRead"`
	CreateAt  time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
	UpdateAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateAt"`
}

// TableName 表名
func (Notification) TableName() string { return "notifications" }

// 通知类型常量
const (
	NotificationTypeCommentPost    = 1 // 评论文章
	NotificationTypeReplyComment   = 2 // 回复评论
	NotificationTypeLikePost       = 3 // 点赞文章
	NotificationTypeLikeComment    = 4 // 点赞评论
)

// 目标类型常量
const (
	NotificationTargetTypePost    = 1 // 文章
	NotificationTargetTypeComment = 2 // 评论
)

// 创建通知请求结构
type CreateNotificationRequest struct {
	UserID     uint   `json:"userId" binding:"required"`
	Type       int    `json:"type" binding:"required,oneof=1 2 3 4"`
	TargetID   uint   `json:"targetId" binding:"required"`
	TargetType int    `json:"targetType" binding:"required,oneof=1 2"`
	Content    string `json:"content" binding:"omitempty,max=500"`
}

// 通知响应结构
type NotificationInfo struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"userId"`
	Type       int    `json:"type"`
	TargetID   uint   `json:"targetId"`
	TargetType int    `json:"targetType"`
	Content    string `json:"content"`
	IsRead     int    `json:"isRead"`
	CreateAt   string `json:"createAt"`
	UpdateAt   string `json:"updateAt"`
}

// 通知列表响应
type NotificationListResponse struct {
	TotalCount    int64             `json:"totalCount"`
	HasMore       bool              `json:"hasMore"`
	Notifications []NotificationInfo `json:"notifications"`
}
