package model

import "time"

// User 用户模型
type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"column:user_id;type:varchar(36);not null;uniqueIndex:idx_users_userID;comment:用户唯一 ID" json:"userID"`
	Username   string    `gorm:"column:username;type:varchar(50);not null;uniqueIndex:idx_users_username;comment:用户名" json:"username"`
	Password   string    `gorm:"column:password;type:varchar(255);not null;comment:密码" json:"-"`
	NickName   string    `gorm:"column:nick_name;type:varchar(50);comment:昵称" json:"nickName"`
	Email      string    `gorm:"column:email;type:varchar(100);comment:邮箱" json:"email"`
	Phone      string    `gorm:"column:phone;type:varchar(20);comment:手机" json:"phone"`
	Role       int       `gorm:"column:role;type:tinyint unsigned;not null;default:1;comment:角色：1-普通用户 2-管理员" json:"role"`
	Status     int       `gorm:"column:status;type:tinyint unsigned;not null;default:0;comment:状态：0-待审核 1-正常 2-封禁" json:"status"`
	GithubOpenID string `gorm:"column:github_openid;type:varchar(50);comment:GitHub OpenID" json:"-"`
	Avatar     string    `gorm:"column:avatar;type:varchar(255);comment:头像 URL" json:"avatar"`
	Bio        string    `gorm:"column:bio;type:varchar(500);comment:个人简介" json:"bio"`
	CreateAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
	UpdateAt   time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateAt"`
}

// TableName 表名
func (User) TableName() string { return "users" }

// 创建用户请求结构
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Nickname string `json:"nick_name" binding:"required,min=2,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty,min=11,max=20"`
}

// 修改用户请求结构
type UpdateUser struct {
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=30"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,min=11,max=20"`
}

// 用户响应结构体
type UserInfo struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     int    `json:"role"`
	Status   int    `json:"status"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	CreateAt string `json:"createAt"`
}

// 管理员审核用户请求
type AuditUserRequest struct {
	UserID uint `json:"userId" binding:"required"`
	Status int  `json:"status" binding:"required,oneof=1 2"` // 1-通过 2-拒绝
}

// 封禁用户请求
type BanUserRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

//  用户登录请求结构

type UserLoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type UserLoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6,max=30"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=30"`
}

// 查询用户列表请求结构体
type ListUserResponse struct {
	TotalCount int64      `json:"totalCount"`
	HasMore    bool       `json:"hasMore"`
	User       []UserInfo `json:"users"`
}
