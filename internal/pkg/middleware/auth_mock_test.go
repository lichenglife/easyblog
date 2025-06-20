package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mock UserBiz
type mockUserBiz struct{}

// ChangePassword implements biz.UserBiz.
func (m *mockUserBiz) ChangePassword(ctx context.Context, userID string, user model.ChangePasswordRequest) error {
	panic("unimplemented")
}

// CreateUser implements biz.UserBiz.
func (m *mockUserBiz) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserInfo, error) {
	panic("unimplemented")
}

// DeleteUser implements biz.UserBiz.
func (m *mockUserBiz) DeleteUser(ctx context.Context, userID string) error {
	panic("unimplemented")
}

// GetUserByID implements biz.UserBiz.
func (m *mockUserBiz) GetUserByID(ctx context.Context, userID string) (*model.UserInfo, error) {
	panic("unimplemented")
}

// GetUserByUsername implements biz.UserBiz.
func (m *mockUserBiz) GetUserByUsername(ctx context.Context, username string) (*model.UserInfo, error) {
	panic("unimplemented")
}

// ListUsers implements biz.UserBiz.
func (m *mockUserBiz) ListUsers(ctx context.Context, page int, pageSize int) (*model.ListUserResponse, error) {
	panic("unimplemented")
}

// UpdateUser implements biz.UserBiz.
func (m *mockUserBiz) UpdateUser(ctx context.Context, user *model.UpdateUser) error {
	panic("unimplemented")
}

// UserLogin implements biz.UserBiz.
func (m *mockUserBiz) UserLogin(ctx context.Context, user model.UserLoginRequest) (*model.UserInfo, error) {
	panic("unimplemented")
}

func (m *mockUserBiz) SomeUserMethod() {}

const (
	testUserID   = "12345"
	testUsername = "testuser"
	testSecret   = "testsecret"
)

func getJWTStrategy() *JWTStrategy {
	// 由于不能直接赋值给JWTSecret，这里通过NewJWTStrategy传递secret
	return NewJWTStrategy(&mockUserBiz{})
}

func TestGenerateTokenAndParseToken(t *testing.T) {
	strategy := getJWTStrategy()
	token, err := strategy.GenerateToken(testUserID, testUsername)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := strategy.ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, testUserID, claims.Subject)
	assert.Equal(t, testUsername, claims.ID)
	assert.Equal(t, "easyblog", claims.Issuer)
	assert.WithinDuration(t, time.Now().Add(time.Hour*2), claims.ExpiresAt.Time, time.Minute)
}

func TestParseToken_InvalidToken(t *testing.T) {
	strategy := getJWTStrategy()
	_, err := strategy.ParseToken("invalid.token.string")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errno.ErrInvalidToken))
}

func TestParseToken_ExpiredToken(t *testing.T) {
	strategy := getJWTStrategy()
	claims := jwt.RegisteredClaims{
		Issuer:    "easyblog",
		Subject:   testUserID,
		ID:        testUsername,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // 已过期
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)

	_, err = strategy.ParseToken(tokenString)
	assert.Error(t, err)
}

func TestAuthMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	strategy := getJWTStrategy()
	token, err := strategy.GenerateToken(testUserID, testUsername)
	require.NoError(t, err)

	r := gin.New()
	r.Use(Auth(strategy))
	r.GET("/test", func(c *gin.Context) {
		username, _ := c.Get("username")
		userID, _ := c.Get("userID")
		c.JSON(200, gin.H{
			"username": username,
			"userID":   userID,
		})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), testUsername)
	assert.Contains(t, w.Body.String(), testUserID)
}

func TestAuthMiddleware_MissingAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	strategy := getJWTStrategy()
	r := gin.New()
	r.Use(Auth(strategy))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthMiddleware_InvalidAuthorizationFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	strategy := getJWTStrategy()
	r := gin.New()
	r.Use(Auth(strategy))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormatToken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	strategy := getJWTStrategy()
	r := gin.New()
	r.Use(Auth(strategy))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.string")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}
