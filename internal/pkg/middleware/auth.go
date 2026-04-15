package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/pkg/auth"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"go.uber.org/zap"
)

// Auth JWT 认证中间件
func Auth(jwt *auth.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Debug("missing authorization header",
				zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errno.ErrUnauthorized.Code(),
				"message": "未提供认证信息",
			})
			c.Abort()
			return
		}

		// 检查 Bearer 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logger.Debug("invalid authorization format",
				zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errno.ErrUnauthorized.Code(),
				"message": "认证信息格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 Token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			logger.Debug("invalid token",
				zap.String("path", c.Request.URL.Path),
				zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errno.ErrUnauthorized.Code(),
				"message": "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RBAC 基于角色的访问控制中间件
func RBAC(requiredRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    errno.ErrPermissionDenied.Code(),
				"message": "未认证用户",
			})
			c.Abort()
			return
		}

		userRole, ok := role.(int)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    errno.ErrPermissionDenied.Code(),
				"message": "角色类型错误",
			})
			c.Abort()
			return
		}

		// 检查角色权限
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				c.Next()
				return
			}
		}

		// 角色权限不足
		c.JSON(http.StatusForbidden, gin.H{
			"code":    errno.ErrPermissionDenied.Code(),
			"message": "权限不足",
		})
		c.Abort()
	}
}

// RateLimit 限流中间件（简单实现）
// 注意：此函数已在 middleware.go 中定义，此处仅为占位符
// 实际使用时请使用 middleware.go 中的 RateLimit 或实现分布式限流
func RateLimitMiddleware(enabled bool, rate, burst int) gin.HandlerFunc {
	// TODO: 使用 redis 实现分布式限流
	// 当前仅作为占位符
	return func(c *gin.Context) {
		if !enabled {
			c.Next()
			return
		}
		// 简单限流逻辑将在后续使用 redis 实现
		c.Next()
	}
}
