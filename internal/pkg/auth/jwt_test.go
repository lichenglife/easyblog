package auth

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestJWT() *JWT {
	config := viper.New()
	config.Set("jwt.secret", "test-secret-key")
	config.Set("jwt.expire", 3600)
	config.Set("jwt.issuer", "easyblog-test")
	return NewJWT(config)
}

func TestJWT_CreateToken_Success(t *testing.T) {
	// Arrange
	jwt := setupTestJWT()

	// Act
	token, err := jwt.CreateToken("user-001", "testuser", 1)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWT_ParseToken_Success(t *testing.T) {
	// Arrange
	jwt := setupTestJWT()
	tokenString, _ := jwt.CreateToken("user-002", "parseuser", 2)

	// Act
	claims, err := jwt.ParseToken(tokenString)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-002", claims.UserID)
	assert.Equal(t, "parseuser", claims.Username)
	assert.Equal(t, 2, claims.Role)
}

func TestJWT_ParseToken_InvalidToken(t *testing.T) {
	// Arrange
	jwt := setupTestJWT()

	// Act
	claims, err := jwt.ParseToken("invalid-token")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWT_ParseToken_ExpiredToken(t *testing.T) {
	// Arrange
	config := viper.New()
	config.Set("jwt.secret", "test-secret-key")
	config.Set("jwt.expire", -1) // 过期时间设置为过去
	config.Set("jwt.issuer", "easyblog-test")
	jwt := &JWT{
		Secret:      config.GetString("jwt.secret"),
		Expire:      config.GetInt64("jwt.expire"),
		Issuer:      config.GetString("jwt.issuer"),
		ExpireTime:  time.Duration(config.GetInt64("jwt.expire")) * time.Second,
	}

	tokenString, _ := jwt.CreateToken("user-003", "expireduser", 1)
	time.Sleep(100 * time.Millisecond) // 等待 token 过期

	// Act
	claims, err := jwt.ParseToken(tokenString)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWT_RefreshToken_Success(t *testing.T) {
	// Arrange
	jwt := setupTestJWT()
	tokenString, _ := jwt.CreateToken("user-004", "refreshuser", 1)

	// Act
	newToken, err := jwt.RefreshToken(tokenString)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)

	// 验证新 token 是否有效
	claims, err := jwt.ParseToken(newToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-004", claims.UserID)
	assert.Equal(t, "refreshuser", claims.Username)
}

func TestJWT_RefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	jwt := setupTestJWT()

	// Act
	newToken, err := jwt.RefreshToken("invalid-token")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, newToken)
}

func TestNewJWT(t *testing.T) {
	// Arrange
	config := viper.New()
	config.Set("jwt.secret", "test-secret")
	config.Set("jwt.expire", 7200)
	config.Set("jwt.issuer", "easyblog")

	// Act
	jwt := NewJWT(config)

	// Assert
	assert.NotNil(t, jwt)
	assert.Equal(t, "test-secret", jwt.Secret)
	assert.Equal(t, int64(7200), jwt.Expire)
	assert.Equal(t, "easyblog", jwt.Issuer)
	assert.Equal(t, 7200*time.Second, jwt.ExpireTime)
}
