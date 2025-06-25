package grpc

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIDInterceptor grpc拦截器，用于设置请求ID
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var requestID string
		md, _ := metadata.FromIncomingContext(ctx)
		// 从请求中解析requestID
		if requestIDs := md["X-Request-ID"]; len(requestIDs) > 0 {
			requestID = requestIDs[0]
		}
		// 生成requestID并设置header
		if requestID == "" {
			requestID = uuid.New().String()
			md.Append("X-Request-ID", requestID)
		}
		// 将元数据设置为新的 incoming context
		ctx = metadata.NewIncomingContext(ctx, md)

		// 设置grpc 请求header
		_ = grpc.SetHeader(ctx, md)
		// 将请求ID 添加到ctx
		//nolint: staticcheck
		ctx = context.WithValue(ctx, "X-Request-ID", requestID)
		// 继续请求
		res, err := handler(ctx, req)
		if err != nil {
			return res, err
		}

		return res, nil
	}

}

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
