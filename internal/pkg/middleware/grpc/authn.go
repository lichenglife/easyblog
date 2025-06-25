package grpc

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"google.golang.org/grpc"
)

var (

	// JWTSecret  TODO 配置
	JWTSecret = "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLi"
)

// UserRetriever 根据用户查询用户信息
type UserRetriever interface {
	// GetUser根据用户ID 获取用户
	GetUser(ctx context.Context, userID string) (*model.User, error)
}

// AuthnInterceptor grpc请求拦截器，认证拦截器
func AuthnInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 解析 JWT token
		token, err := auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return nil, fmt.Errorf("获取token失败 %v", err)
		}
		claims, err := ParseToken(token)
		if err != nil {
			return nil, fmt.Errorf("解析token失败 %v", err)
		}
		userID := claims.Subject
		username := claims.ID
		// user, err := retriever.GetUser(ctx, userID)
		// if err != nil {
		// 	return nil, errno.ErrUnauthorized
		// }
		// 设置用户ID到上下文中
		//nolint:staticcheck // SA1029: string key for context is required for Gin compatibility
		ctx = context.WithValue(ctx, "userID", userID)
		//nolint:staticcheck // SA1029: string key for context is required for Gin compatibility
		ctx = context.WithValue(ctx, "username", username)
		// 继续处理请求
		return handler(ctx, req)
	}
}

// ParseToken  解析token

func ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	// 解析token
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, errno.ErrInvalidToken
	}
	// 校验token是否有效
	if !token.Valid {
		return nil, errno.ErrInvalidToken
	}

	// 提取声明
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
