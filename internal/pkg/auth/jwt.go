package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Claims JWT 声明
type Claims struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

// JWT JWT 结构
type JWT struct {
	Secret   string
	Expire   int64
	Issuer   string
	ExpireTime time.Duration
}

// NewJWT 创建 JWT 实例
func NewJWT(config *viper.Viper) *JWT {
	return &JWT{
		Secret:      config.GetString("jwt.secret"),
		Expire:      config.GetInt64("jwt.expire"),
		Issuer:      config.GetString("jwt.issuer"),
		ExpireTime:  time.Duration(config.GetInt64("jwt.expire")) * time.Second,
	}
}

// CreateToken 创建 Token
func (j *JWT) CreateToken(userID, username string, role int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}

// ParseToken 解析 Token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return j.CreateToken(claims.UserID, claims.Username, claims.Role)
}
