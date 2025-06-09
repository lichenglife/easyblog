package validation

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Validator 全局校验器
var validate *validator.Validate

func init() {
	// 从 Gin 的 binding.Validator 获取 validator 实例
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v

		// 注册自定义校验规则
		_ = validate.RegisterValidation("username", validateUsername)
		_ = validate.RegisterValidation("phone", validatePhone)
		_ = validate.RegisterValidation("email", validateEmail)
		_ = validate.RegisterValidation("password", validatePassword)
	} else {
		panic("Failed to initialize validator")
	}
}

// ValidateStruct 校验结构体
func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}

// ValidationErrors 格式化校验错误
func ValidationErrors(err error) map[string]string {
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = translateError(fieldErr)
		}
		return errors
	}

	errors["error"] = err.Error()
	return errors
}

// translateError 转换校验错误为友好提示
func translateError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "字段是必填项"
	case "username":
		return "用户名格式不正确"
	case "phone":
		return "手机号格式不正确"
	case "email":
		return "邮箱格式不正确"
	case "password":
		return "密码格式不正确"
	case "min":
		return fmt.Sprintf("长度不能小于 %s", err.Param())
	case "max":
		return fmt.Sprintf("长度不能大于 %s", err.Param())
	default:
		return "格式不正确"
	}
}

// validateUsername 校验用户名格式
func validateUsername(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return re.MatchString(fl.Field().String())
}

// validatePhone 校验手机号格式
func validatePhone(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(fl.Field().String())
}

// validateEmail 校验邮箱格式
func validateEmail(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(fl.Field().String())
}

// validatePassword 校验密码格式
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 6 || len(password) > 30 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
