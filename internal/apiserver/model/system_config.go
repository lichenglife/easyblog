package model

import "time"

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ConfigKey  string    `gorm:"column:config_key;type:varchar(50);not null;uniqueIndex:idx_config_key;comment:配置键" json:"configKey"`
	ConfigValue string   `gorm:"column:config_value;type:text;not null;comment:配置值" json:"configValue"`
	ConfigType int       `gorm:"column:config_type;type:tinyint unsigned;not null;default:1;comment:配置类型：1-字符串 2-数字 3-布尔 4-JSON" json:"configType"`
	GroupName  string    `gorm:"column:group_name;type:varchar(50);comment:分组名称" json:"groupName"`
	Description string   `gorm:"column:description;type:varchar(255);comment:配置描述" json:"description"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"`
	UpdatedAtBy uint      `gorm:"column:updated_at_by;type:bigint unsigned;comment:最后更新人 ID" json:"updatedBy"`
}

// TableName 表名
func (SystemConfig) TableName() string { return "system_configs" }

// 系统配置请求结构
type UpdateSystemConfigRequest struct {
	ConfigKey   string `json:"configKey" binding:"required"`
	ConfigValue string `json:"configValue" binding:"required"`
}

// 系统配置响应结构
type SystemConfigInfo struct {
	ID          uint      `json:"id"`
	ConfigKey   string    `json:"configKey"`
	ConfigValue string    `json:"configValue"`
	ConfigType  int       `json:"configType"`
	GroupName   string    `json:"groupName"`
	Description string    `json:"description"`
	UpdatedAt   string    `json:"updatedAt"`
	UpdatedBy   uint      `json:"updatedBy"`
}

// 系统配置列表响应
type SystemConfigListResponse struct {
	TotalCount int64              `json:"totalCount"`
	Configs    []SystemConfigInfo `json:"configs"`
}
