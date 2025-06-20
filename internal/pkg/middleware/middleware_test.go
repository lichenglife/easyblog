package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {

	//mock request
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Recovery(log.Log)).Use(RequestID()).Use(Logger(log.Log))
	r.GET("/test", func(c *gin.Context) {
		username, _ := c.Get("username")
		userID, _ := c.Get("userID")
		c.JSON(200, gin.H{
			"username": username,
			"userID":   userID,
		})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
