package core

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 定义参数正则表达式
var (
	// 手机号正则表达式
	phoneRegexp = regexp.MustCompile(`^1[3-9]\d{9}$`)
	// 邮箱正则表达式
	emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// 密码正则表达式（8-20位字母数字）
	passwordRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{8,20}$`)
	// 用户名正则表达式 4-16位 包含大小写字母、数字、下划线
	usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]{4,16}$`)
)

// 初始化校验器
func InitValidator() {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phone", validatePhone)
		_ = v.RegisterValidation("email", validateEmail)
		_ = v.RegisterValidation("password", validatePassword)
		_ = v.RegisterValidation("username", validateUsername)

		_ = v.RegisterValidation("postTitle", validatePostTitle)
		_ = v.RegisterValidation("postContent", validatePostContent)
	}

}

// validatePhone 校验手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return phoneRegexp.MatchString(phone)
}

// validateEmail 校验邮箱
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return emailRegexp.MatchString(email)
}

// validatePassword 校验密码
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return passwordRegexp.MatchString(password)
}

// validateUsername 校验用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return usernameRegexp.MatchString(username)
}

// validatePostTitle 校验帖子标题
func validatePostTitle(fl validator.FieldLevel) bool {
	title := fl.Field().String()
	return len(title) >= 1 && len(title) <= 100
}

// validatePostContent 校验帖子内容
func validatePostContent(fl validator.FieldLevel) bool {
	content := fl.Field().String()
	return len(content) >= 1 && len(content) <= 1000
}
