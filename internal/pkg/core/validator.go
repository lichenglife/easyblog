package core

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
)

// 定义参数正则表达式
var (
	// 手机号正则表达式
	phoneRegexp = regexp.MustCompile(`^1[3-9]\d{9}$`)
	// 邮箱正则表达式
	emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// 用户名正则表达式 4-16位 包含大小写字母、数字、下划线
	usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]{4,16}$`)
	// 密码正则表达式（8-20位字母数字）
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9]{8,20}$`)
	// 小写字母检查
	lowercaseRegex = regexp.MustCompile(`[a-z]`)
	// 大写字母检查
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)
	// 数字检查
	digitRegex = regexp.MustCompile(`[0-9]`)

	// 全局验证器实例
	validate *validator.Validate
)

// 初始化校验器
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		_ = validate.RegisterValidation("phone", validatePhone)
		_ = validate.RegisterValidation("email", validateEmail)
		_ = validate.RegisterValidation("password", validatePassword)
		_ = validate.RegisterValidation("username", validateUsername)
		_ = validate.RegisterValidation("postTitle", validatePostTitle)
		_ = validate.RegisterValidation("postContent", validatePostContent)
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
	return passwordRegex.MatchString(password)
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

// BindAndValid 绑定并验证请求参数
func BindAndValid(c *gin.Context, obj interface{}) error {
	// 根据Content-Type选择合适的绑定器
	if err := c.ShouldBind(obj); err != nil {
		// 处理验证错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, validationError := range validationErrors {
				field := validationError.Field()
				tag := validationError.Tag()
				param := validationError.Param()
				value := validationError.Value()

				switch tag {
				case "required":
					errMsgs = append(errMsgs, formatFieldError(field, "不能为空"))
				case "min":
					errMsgs = append(errMsgs, formatFieldError(field, "长度不能小于 "+param+" 个字符"))
				case "max":
					errMsgs = append(errMsgs, formatFieldError(field, "长度不能大于 "+param+" 个字符"))
				case "phone":
					errMsgs = append(errMsgs, formatFieldError(field, "格式不正确，必须是11位有效的手机号码"))
				case "username":
					errMsgs = append(errMsgs, formatFieldError(field, "格式不正确，只能包含字母、数字和下划线，长度4-16位"))
				case "password":
					var details []string
					if value != nil {
						strValue, ok := value.(string)
						if ok && len(strValue) > 0 {
							if !lowercaseRegex.MatchString(strValue) {
								details = append(details, "必须包含小写字母")
							}
							if !uppercaseRegex.MatchString(strValue) {
								details = append(details, "必须包含大写字母")
							}
							if !digitRegex.MatchString(strValue) {
								details = append(details, "必须包含数字")
							}
							if !passwordRegex.MatchString(strValue) {
								details = append(details, "长度必须在8-20位之间")
							}
						}
					}
					if len(details) > 0 {
						errMsgs = append(errMsgs, formatFieldError(field, "格式不正确："+strings.Join(details, "，")))
					} else {
						errMsgs = append(errMsgs, formatFieldError(field, "格式不正确，必须包含大小写字母和数字，长度8-20位"))
					}
				case "email":
					errMsgs = append(errMsgs, formatFieldError(field, "格式不正确，请输入有效的电子邮箱地址"))
				default:
					errMsgs = append(errMsgs, formatFieldError(field, "验证失败"))
				}
			}
			if len(errMsgs) > 0 {
				return errno.ErrInvalidParams.WithMessage("参数验证失败：" + strings.Join(errMsgs, "；"))
			}
		}
		return errno.ErrInvalidParams.WithMessage("参数验证失败：" + err.Error())
	}
	return nil
}

// formatFieldError 格式化字段错误信息
func formatFieldError(field, msg string) string {
	switch field {
	case "Username":
		return "用户名" + msg
	case "Password":
		return "密码" + msg
	case "Email":
		return "邮箱" + msg
	case "Phone":
		return "手机号" + msg
	case "Nickname":
		return "昵称" + msg
	default:
		return field + msg
	}
}
