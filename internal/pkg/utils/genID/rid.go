package genid

import "github.com/google/uuid"

//  生成UID 以及POSTID  工具类

type ResourceID string

const (
	UserIDPrefix ResourceID = "user" // 用户ID前缀
	PostIDPrefix ResourceID = "post" // 帖子ID前缀
)

func (r ResourceID) String() string {
	return string(r)
}

// GenerateUserID 生成用户ID
// 生成格式为 "user-<unique_id>"，其中 <unique_id> 是一个唯一标识符
func generateUniqueID() string {
	// 使用uuid 生成随机6位数字符串
	return uuid.New().String()[:6]
}

func GenerateUserID() string {
	return string(UserIDPrefix) + "-" + generateUniqueID()
}

func GeneratePostID() string {
	return string(PostIDPrefix) + "-" + generateUniqueID()
}
