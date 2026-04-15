package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRequestID_NewRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		requestID := c.GetString("requestID")
		c.JSON(http.StatusOK, gin.H{"requestID": requestID})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	assert.Contains(t, w.Body.String(), "requestID")
}

func TestRequestID_ExistingRequestID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/test", func(c *gin.Context) {
		requestID := c.GetString("requestID")
		c.JSON(http.StatusOK, gin.H{"requestID": requestID})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "custom-request-id-123")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "custom-request-id-123", w.Header().Get("X-Request-ID"))
	assert.Contains(t, w.Body.String(), "custom-request-id-123")
}

func TestLogger_RequestLogging(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)
	router.Use(Logger(logger))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?query=value", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecovery_NoPanic(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)
	router.Use(Recovery(logger))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestRecovery_WithPanic(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)
	router.Use(Recovery(logger))
	router.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "服务器内部错误")
}

func TestRecovery_ValidationError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)
	router.Use(Recovery(logger))
	router.GET("/test", func(c *gin.Context) {
		// 模拟 validator.ValidationErrors 类型的错误
		err := validator.ValidationErrors{}
		panic(&err)
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "请求参数验证失败")
}

func TestRecovery_ErrorType(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)
	router.Use(Recovery(logger))
	router.GET("/test", func(c *gin.Context) {
		panic(errors.New("custom error message"))
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "custom error message")
}

func TestCORS_AllowAllOrigins(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_OPTIONSRequest(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.OPTIONS("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_AllowHeaders(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
	assert.Contains(t, allowHeaders, "Authorization")
	assert.Contains(t, allowHeaders, "Content-Type")
}

func TestRateLimit_PassThrough(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimit(100, 200))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_Serial(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)

	router.Use(RequestID())
	router.Use(Logger(logger))
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		requestID := c.GetString("requestID")
		c.JSON(http.StatusOK, gin.H{
			"requestID": requestID,
			"status":    "ok",
		})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestLogger_QueryParams(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)
	router.Use(Logger(logger))
	router.GET("/test", func(c *gin.Context) {
		query := c.Query("q")
		c.JSON(http.StatusOK, gin.H{"query": query})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?q=search", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "search")
}

func TestLogger_ClientIP(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)
	router.Use(Logger(logger))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecovery_HTTPMethod(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)
	router.Use(Recovery(logger))

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	for _, method := range methods {
		router.Handle(method, "/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"method": method})
		})
	}

	// Act & Assert for each method
	for _, method := range methods {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestCORS_AllowMethods(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	// Assert
	allowMethods := w.Header().Get("Access-Control-Allow-Methods")
	assert.True(t, strings.Contains(allowMethods, "POST"))
	assert.True(t, strings.Contains(allowMethods, "GET"))
	assert.True(t, strings.Contains(allowMethods, "PUT"))
	assert.True(t, strings.Contains(allowMethods, "DELETE"))
	assert.True(t, strings.Contains(allowMethods, "OPTIONS"))
}
