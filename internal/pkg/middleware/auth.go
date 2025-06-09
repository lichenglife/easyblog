package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	biz "github.com/lichenglife/easyblog/internal/apiserver/biz/v1/user"
	"github.com/lichenglife/easyblog/internal/pkg/core"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
)

const (

	// JWTSecret  TODO 配置
	JWTSecret = "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLi"
	// 请求头Authorization
	Authorization = "Authorization"

	AuthorizationBearer = "Bearer"
)

// 通过gin中间件实现jwt token的生成、解析和验证

type AuthStrategy interface {
	// GenerateToken 生成JWT Token
	GenerateToken(userID, username string) (string, error)
	// ParseToken 解析JWT Token
	ParseToken(tokenString string) (*jwt.RegisteredClaims, error)
	// ValidateToken 验证JWT Token
	//ValidateToken(tokenString string) (bool, error)
}

// Auth 认证中间件
func Auth(strategy AuthStrategy) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 获取请求头header
		authorization := c.GetHeader(Authorization)
		if authorization == "" {
			core.WriteResponse(c, errno.ErrUnauthorized, nil)
			c.Abort()
			return
		}
		// 解析请求头中的token
		splits := strings.Split(authorization, " ")
		if !(len(splits) == 2 && splits[0] == AuthorizationBearer) {
			core.WriteResponse(c, errno.ErrUnauthorized, nil)
			c.Abort()
			return
		}
		claims, err := strategy.ParseToken(splits[1])
		if err != nil {
			core.WriteResponse(c, errno.ErrUnauthorized, nil)
			c.Abort()
			return
		}
		// 设置用户ID到上下文中
		c.Set("username", claims.ID)
		c.Set("userID", claims.Subject)
		c.Set("user", claims.Subject)
		c.Next()

	}
}

type JWTStrategy struct {
	// 实现AuthStrategy接口的方法
	userbiz biz.UserBiz // 假设有一个UserBiz接口用于用户业务逻辑
}

// NewJWTStrategy 创建一个新的JWT策略实例
func NewJWTStrategy(userbiz biz.UserBiz) *JWTStrategy {
	return &JWTStrategy{userbiz: userbiz}
}

// GenerateToken 生成JWT Token
func (j *JWTStrategy) GenerateToken(userID, username string) (string, error) {
	// 实现生成JWT Token的逻辑
	// 生成jwtClaims
	claims := jwt.RegisteredClaims{
		// 签发者
		Issuer: "easyblog",
		// 主体
		Subject: userID,
		// ID
		ID: username,
		// 到期时间
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
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
		return "", errno.ErrInvalidToken
	}

	return tokenString, nil
}

// ParseToken  解析token

func (j *JWTStrategy) ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {
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
