package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// 生成 UserID
func GenUserID() string {
	return fmt.Sprintf("user-%s", GenShortID())
}

// 生成PostID
func GenPostID() string {
	return fmt.Sprintf("post-%s", GenShortID())
}

// GenShortID 生成短ID
func GenShortID() string {
	id := uuid.New()
	// 获取UUID的前8位
	return id.String()[:8]
}
