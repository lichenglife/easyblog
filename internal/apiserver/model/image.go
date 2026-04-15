package model

import "time"

// Image 图片模型
type Image struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"column:user_id;type:bigint unsigned;not null;index:idx_images_user_id;comment:上传者 ID" json:"userId"`
	FilePath   string    `gorm:"column:file_path;type:varchar(255);not null;comment:文件路径" json:"filePath"`
	FileName   string    `gorm:"column:file_name;type:varchar(100);not null;comment:文件名" json:"fileName"`
	FileSize   uint      `gorm:"column:file_size;type:int unsigned;not null;comment:文件大小（字节）" json:"fileSize"`
	FileType   string    `gorm:"column:file_type;type:varchar(50);not null;comment:文件类型（MIME）" json:"fileType"`
	Width      uint      `gorm:"column:width;type:int unsigned;comment:图片宽度" json:"width"`
	Height     uint      `gorm:"column:height;type:int unsigned;comment:图片高度" json:"height"`
	Storage    int       `gorm:"column:storage;type:tinyint unsigned;not null;default:1;comment:存储方式：1-本地 2-OSS" json:"storage"`
	URL        string    `gorm:"column:url;type:varchar(500);not null;comment:访问 URL" json:"url"`
	Status     int       `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:状态：0-禁用 1-正常" json:"status"`
	CreateAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
}

// TableName 表名
func (Image) TableName() string { return "images" }

// 上传图片响应结构
type ImageInfo struct {
	ID       uint      `json:"id"`
	UserID   uint      `json:"userId"`
	FilePath string    `json:"filePath"`
	FileName string    `json:"fileName"`
	FileSize uint      `json:"fileSize"`
	FileType string    `json:"fileType"`
	Width    uint      `json:"width"`
	Height   uint      `json:"height"`
	Storage  int       `json:"storage"`
	URL      string    `json:"url"`
	CreateAt string    `json:"createAt"`
}

// 图片列表响应
type ImageListResponse struct {
	TotalCount int64      `json:"totalCount"`
	HasMore    bool       `json:"hasMore"`
	Images     []ImageInfo `json:"images"`
}
