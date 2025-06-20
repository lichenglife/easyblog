package middleware

//  单元测试的原则和方法
// 1. 测试的是函数的行为，而不是函数的实现
// 2. 测试的是函数的输入和输出，而不是函数的内部实现
// 3. 测试的是函数的边界情况，而不是函数的内部实现
// 4. 测试的是函数的错误处理，而不是函数的内部实现
// 如何进行单元测试代码中的依赖、成员进行mock,分场景进行介绍，函数、外部依赖、全局变量等如何分别进行mock

// 单元测试分类
// 1. 函数的测试
// 2. 外部依赖的测试
// 3. 全局变量的测试
// 4. 边界情况的测试
// 5. 错误处理的测试
// 6. 性能测试
// 7. 代码覆盖率测试
// 8. 集成测试
// 9. 基准测试
// 10. 文档测试
// 11. 重构测试
// 12. 并发测试
// 13. 测试工具的选择

// -----------------------------
// GoMock + mockgen 单元测试实践
// -----------------------------
//
// 步骤1：用 mockgen 生成接口 mock 文件
// 例如：
// mockgen -source=internal/apiserver/biz/v1/user/user.go -destination=internal/apiserver/biz/v1/user/mocks/mock_userbiz.go -package=mocks
//
// 步骤2：在测试文件中引入 mock 包和 gomock
import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/lichenglife/easyblog/internal/apiserver/biz/v1/user/mocks"
	"github.com/lichenglife/easyblog/internal/pkg/errno"
	"github.com/stretchr/testify/assert"
)

func setJWTStrategy(t *testing.T) *JWTStrategy {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBiz := mocks.NewMockUserBiz(ctrl)
	JWTStrategy := NewJWTStrategy(mockBiz)
	return JWTStrategy
}
func TestNewJWTStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBiz := mocks.NewMockUserBiz(ctrl)
	NewJWTStrategy(mockBiz)
}

// 步骤3：在测试函数中创建 mock 对象
func TestGenerateTokenAndParseTokenByMockUserBiz(t *testing.T) {
	strategy := setJWTStrategy(t)
	token, err := strategy.GenerateToken("u1", "testuser")
	assert.NoError(t, err)

	// paras
	claims, err := strategy.ParseToken(token)
	t.Log(claims)
	assert.NoError(t, err)
	assert.Equal(t, claims.ID, "testuser")

}

func TestParseToken_InvalidToken_ByMockUserBiz(t *testing.T) {
	strategy := setJWTStrategy(t)
	_, err := strategy.ParseToken("invalid.token")
	assert.Error(t, err)
	assert.EqualError(t, err, errno.ErrInvalidToken.Error())

}
func TestParseToken_ExpiredToken_ByMockUserBiz(t *testing.T) {
	//   token过期的场景
	strategy := setJWTStrategy(t)
	claims := jwt.RegisteredClaims{
		// 签发者
		Issuer: "easyblog",
		// 主体
		Subject: "U1",
		// ID
		ID: "U1",
		// 到期时间
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * -2)),
		// 签发时间
		IssuedAt: jwt.NewNumericDate(time.Now().Add(time.Hour * -2)),
		//
		NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * -2)),
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 创建token签名
	tokenString, _ := token.SignedString([]byte(JWTSecret))
	_, err := strategy.ParseToken(tokenString)
	assert.Error(t, err)

}

func TestAuthMiddleware_SUCCESS(t *testing.T) {

	gin.SetMode(gin.TestMode)

	strategy := setJWTStrategy(t)
	tokenString, err := strategy.GenerateToken("userID-1", "zhangsan")
	assert.NoError(t, err)

	// mock request
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
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	r.ServeHTTP(w, req)

	// 设置断言
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "zhangsan")
	assert.Contains(t, w.Body.String(), "userID-1")

}

func Test_AuthMiddleware_WithoutToken(t *testing.T) {

	gin.SetMode(gin.TestMode)

	strategy := setJWTStrategy(t)
	//tokenString, err := strategy.GenerateToken("userID-1", "zhangsan")
	//assert.NoError(t, err)

	// mock request
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
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer ")
	r.ServeHTTP(w, req)

	// 设置断言
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func Test_AuthMiddleware_InvalidToken(t *testing.T) {

	gin.SetMode(gin.TestMode)

	strategy := setJWTStrategy(t)
	tokenString, err := strategy.GenerateToken("userID-1", "zhangsan")
	assert.NoError(t, err)

	// mock request
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
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString+"test")
	r.ServeHTTP(w, req)

	// 设置断言
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func Test_AuthMiddleware_InvalidHeader(t *testing.T) {

	gin.SetMode(gin.TestMode)

	strategy := setJWTStrategy(t)

	// mock request
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
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	r.ServeHTTP(w, req)

	// 设置断言
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}
