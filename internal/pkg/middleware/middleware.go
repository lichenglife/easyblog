package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/lichenglife/easyblog/internal/pkg/log"
)

// RequestID 生成请求ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger 记录请求日志
func Logger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// Recovery 恢复panic
func Recovery(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("requestID", c.GetString("requestID")),
				)

				// 根据错误类型返回友好的错误信息
				var errMsg string
				switch e := err.(type) {
				case *validator.ValidationErrors:
					errMsg = "请求参数验证失败"
				case error:
					errMsg = e.Error()
				default:
					errMsg = "服务器内部错误"
				}

				c.JSON(http.StatusOK, gin.H{
					"code":    errno.ErrInvalidParams.Code(),
					"message": errMsg,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现限流逻辑
		c.Next()
	}
}
