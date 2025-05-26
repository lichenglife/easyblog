package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
)

const (

	// JWTSecret JWT秘钥
	JWTSecret = "easyblog-jwt-secret"
	// 请求头header authorization
	AuthorizationHeader = "Authorization"

	AuthorizationBearer = "Bearer"
)

// AuthStrategy 定义了认证策略接口
type AuthStrategy interface {
	// GenerateToken 生成token
	GenerateToken(userID, username string) (string, error)
	// Parse 解析token
	Parse(tokenString string) (*jwt.RegisteredClaims, error)
}

// Auth  认证中间件
func Auth(authStrategy AuthStrategy) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		// 获取请求头header
		authorization := ctx.GetHeader(AuthorizationHeader)
		if authorization == "" {
			core.WriteResponse(ctx, errno.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
		// 分割请求头
		parts := strings.SplitN(authorization, " ", 2)
		if !(len(parts) == 2 && parts[0] == AuthorizationBearer) {
			core.WriteResponse(ctx, errno.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
		// 解析token

		claims, err := authStrategy.Parse(parts[1])
		if err != nil {
			core.WriteResponse(ctx, errno.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
		// 设置用户ID以及用户到context
		ctx.Set("userID", claims.Subject)
		ctx.Set("username", claims.ID)
		ctx.Set("user", claims.Subject)
		ctx.Next()
	}

}

// JWTStrategy 授权策略
type JWTStrategy struct {
	userBiz biz.UserBiz
}

// NewJWTStrategy  生成授权策略
func NewJWTStrategy(userBiz biz.UserBiz) AuthStrategy {
	return &JWTStrategy{
		userBiz: userBiz,
	}

}

// GenerateToken 生成token
func (js *JWTStrategy) GenerateToken(userID, username string) (string, error) {

	// 生成jwtclaims
	claims := jwt.RegisteredClaims{
		// 签发者
		Issuer: "easyblog",
		// 主体
		Subject: userID,
		// ID
		ID: username,
		// 到期时间
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		// 签发时间
		IssuedAt: jwt.NewNumericDate(time.Now()),
		//
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 创建token签名
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Parse 解析token
func (s *JWTStrategy) Parse(tokenString string) (*jwt.RegisteredClaims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// 检查Token是否有效
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// 提取声明
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
