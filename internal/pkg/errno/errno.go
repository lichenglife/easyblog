package errno

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// Errno 定义了错误码接口
type Errno interface {
	// Error 返回错误信息
	Error() string
	// Code 返回错误码
	Code() int
	// Message 返回错误详情
	Message() string
	// HTTP 返回HTTP状态码
	HTTP() int
	// WithMessage 设置自定义错误消息
	WithMessage(message string) Errno
}

// errno 实现了Errno接口
type errno struct {
	code    int
	message string
	http    int
}

// Error 返回错误信息
func (e *errno) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.code, e.message)
}

// Code 返回错误码
func (e *errno) Code() int {
	return e.code
}

// Message 返回错误详情
func (e *errno) Message() string {
	return e.message
}

// HTTP 返回HTTP状态码
func (e *errno) HTTP() int {
	return e.http
}

// WithMessage 设置自定义错误消息
func (e *errno) WithMessage(message string) Errno {
	return &errno{
		code:    e.code,
		message: message,
		http:    e.http,
	}
}

// New 创建一个新的错误码
func New(code int, message string, http int) Errno {
	return &errno{
		code:    code,
		message: message,
		http:    http,
	}
}

// 定义系统级错误码
var (
	// OK 表示成功
	OK = New(0, "OK", http.StatusOK)
	// ErrInternalServer 表示服务器内部错误
	ErrInternalServer = New(10001, "服务器内部错误", http.StatusInternalServerError)
	// ErrDatabase 表示数据库错误
	ErrDatabase = New(10002, "数据库错误", http.StatusInternalServerError)
	// ErrParam 表示参数错误
	ErrInvalidParams = New(10003, "参数错误", http.StatusBadRequest)
	// ErrUnauthorized 表示未授权
	ErrUnauthorized = New(10004, "未授权", http.StatusUnauthorized)
	// ErrForbidden 表示禁止访问
	ErrForbidden = New(10005, "禁止访问", http.StatusForbidden)
	// ErrNotFound 表示资源不存在
	ErrNotFound = New(10006, "资源不存在", http.StatusNotFound)
	// ErrTooManyRequests 表示请求过于频繁
	ErrTooManyRequests = New(10007, "请求过于频繁", http.StatusTooManyRequests)
	// ErrInvalidToken 表示无效的Token
	ErrInvalidToken = New(10008, "无效的Token", http.StatusUnauthorized)
	// ErrTokenExpired 表示Token已过期
	ErrTokenExpired    = New(10009, "Token已过期", http.StatusUnauthorized)
	ErrEncryptPassword = New(10010, "用户密码加密失败", http.StatusUnauthorized)

	ErrBind    = New(100010, "参数错误", http.StatusBadGateway)
	ErrUnknown = New(99999, "未知错误", http.StatusBadRequest)
)

// 定义业务级错误码
var (
	// 用户相关错误码 (2xxxx)
	ErrUserNotFound      = New(20001, "用户不存在", http.StatusNotFound)
	ErrUserAlreadyExist  = New(20002, "用户已存在", http.StatusConflict)
	ErrPasswordIncorrect = New(20003, "密码错误", http.StatusUnauthorized)
	ErrInvalidUsername   = New(20004, "用户名格式不正确", http.StatusBadRequest)
	ErrInvalidPassword   = New(20005, "密码格式不正确", http.StatusBadRequest)
	ErrInvalidPhone      = New(20006, "手机号格式不正确", http.StatusBadRequest)
	ErrInvalidEmail      = New(20007, "邮箱格式不正确", http.StatusBadRequest)
	ErrGenerateToken     = New(20008, "生成token失败", http.StatusUnauthorized)

	// 博客相关错误码 (3xxxx)
	ErrPostNotFound       = New(30001, "博客不存在", http.StatusNotFound)
	ErrPostAccessDenied   = New(30002, "无权访问该博客", http.StatusForbidden)
	ErrInvalidPostTitle   = New(30003, "博客标题格式不正确", http.StatusBadRequest)
	ErrInvalidPostContent = New(30004, "博客内容格式不正确", http.StatusBadRequest)
)

// IsRecordNotFound 判断是否是记录不存在错误
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// Decode 解码错误
func Decode(err error) Errno {
	if err == nil {
		return OK
	}

	// 尝试转换为Errno类型
	if e, ok := err.(Errno); ok {
		return e
	}

	// 检查是否是记录不存在错误
	if IsRecordNotFound(err) {
		return ErrNotFound
	}

	// 默认返回内部错误
	return ErrInternalServer.WithMessage(err.Error())
}
