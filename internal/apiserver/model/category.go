package model

import "time"

// Category 分类模型
type Category struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"column:name;type:varchar(50);not null;uniqueIndex:idx_categories_name;comment:分类名称" json:"name"`
	Slug     string    `gorm:"column:slug;type:varchar(50);not null;uniqueIndex:idx_categories_slug;comment:分类别名" json:"slug"`
	ParentID uint      `gorm:"column:parent_id;type:bigint unsigned;default:0;index:idx_categories_parent_id;comment:父分类 ID" json:"parentId"`
	Level    int       `gorm:"column:level;type:tinyint unsigned;not null;default:1;comment:层级（最多 3 级）" json:"level"`
	Sort     int       `gorm:"column:sort;type:int unsigned;not null;default:0;comment:排序（越小越靠前）" json:"sort"`
	PostCount uint    `gorm:"column:post_count;type:int unsigned;not null;default:0;comment:文章数" json:"postCount"`
	Status   int       `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:状态：0-禁用 1-启用" json:"status"`
	CreateAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
	UpdateAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateAt"`
}

// TableName 表名
func (Category) TableName() string { return "categories" }

// 创建分类请求结构
type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	Slug     string `json:"slug" binding:"required,max=50"`
	ParentID uint   `json:"parentId" binding:"omitempty"`
	Sort     int    `json:"sort" binding:"omitempty"`
}

// 更新分类请求结构
type UpdateCategoryRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required,max=50"`
	Slug     string `json:"slug" binding:"required,max=50"`
	ParentID uint   `json:"parentId" binding:"omitempty"`
	Sort     int    `json:"sort" binding:"omitempty"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// 分类响应结构
type CategoryInfo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ParentID  uint      `json:"parentId"`
	Level     int       `json:"level"`
	Sort      int       `json:"sort"`
	PostCount uint      `json:"postCount"`
	Status    int       `json:"status"`
	CreateAt  string    `json:"createAt"`
	UpdateAt  string    `json:"updateAt"`
	Children  []CategoryInfo `json:"children,omitempty"`
}

// 分类树响应
type CategoryTreeResponse struct {
	TotalCount int64          `json:"totalCount"`
	Categories []CategoryInfo `json:"categories"`
}
