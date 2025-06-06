package model

import "time"

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"column:userID;type:varchar(36);not null;uniqueIndex:user.userID;comment:用户唯一 ID" json:"userID"`
	Username  string    `gorm:"column:username;type:varchar(36);not null;uniqueIndex:user.username;comment:用户名" json:"username"`
	Password  string    `gorm:"column:password;type:varchar(36);not null;comment:密码" json:"-"`
	NickName  string    `gorm:"column:nickName;type:varchar(36);not null;comment:昵称" json:"nickName"`
	Email     string    `gorm:"column:email;type:varchar(36);not null;comment:邮箱" json:"email"`
	Phone     string    `gorm:"column:phone;type:varchar(36);not null;uniqueIndex:user.phone;comment:手机" json:"phone"`
	CreatedAt time.Time `gorm:"column:createdAt;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updateAt"`
}

// TableName 表名
func (User) TableName() string { return "user" }

// 创建用户请求结构
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,password"`
	Nickname string `json:"nickname" binding:"required,min=2,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,phone"`
}

// 修改用户请求结构
type UpdateUser struct {
	UserID   string `json:"userID" binding:"required,uuid"`
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=30"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,phone"`
}

// 用户响应结构体
type UserInfo struct {
	UserID    string    `json:"userID"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	BlogTotal int       `json:"blogTotal"`
}

//  用户登录请求结构

type UserLoginRequest struct {
	Username string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,min=6,max=30"`
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
