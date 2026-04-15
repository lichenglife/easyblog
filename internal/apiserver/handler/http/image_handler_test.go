package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testImageHandler struct {
	logger    *log.Logger
	uploadDir string
	maxSize   int64
}

func (h *testImageHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	if file.Size > h.maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件超出最大限制"})
		return
	}
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}
	filename := "test_" + filepath.Base(file.Filename)
	savePath := filepath.Join(h.uploadDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"filename": filename, "url": "/images/" + filename, "size": file.Size})
}

func (h *testImageHandler) GetImage(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join(h.uploadDir, filename)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}
	c.Header("Content-Type", "image/png")
	c.File(filepath)
}

func (h *testImageHandler) ListImages(c *gin.Context) {
	entries, _ := os.ReadDir(h.uploadDir)
	var images []gin.H
	for _, entry := range entries {
		if !entry.IsDir() {
			images = append(images, gin.H{"filename": entry.Name(), "url": "/images/" + entry.Name()})
		}
	}
	c.JSON(http.StatusOK, gin.H{"total": len(images), "images": images})
}

func (h *testImageHandler) DeleteImage(c *gin.Context) {
	filename := c.Param("filename")
	filepath := filepath.Join(h.uploadDir, filename)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}
	os.Remove(filepath)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func setupTestImageHandler(t *testing.T) (*gin.Engine, *testImageHandler, string) {
	gin.SetMode(gin.TestMode)
	tempDir := t.TempDir()

	config := viper.New()
	config.Set("log.level", "debug")
	config.Set("log.dir", "/tmp/logs")
	logger, _ := log.NewLogger(config)

	handler := &testImageHandler{
		logger:    logger,
		uploadDir: tempDir,
		maxSize:   5 * 1024 * 1024,
	}

	router := gin.New()
	router.POST("/upload", handler.UploadImage)
	router.GET("/images/:filename", handler.GetImage)
	router.GET("/images", handler.ListImages)
	router.DELETE("/images/:filename", handler.DeleteImage)

	return router, handler, tempDir
}

func TestImageHandler_ListImages_Empty(t *testing.T) {
	router, _, _ := setupTestImageHandler(t)

	req, _ := http.NewRequest(http.MethodGet, "/images", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "images")
}

func TestImageHandler_GetImage_NotFound(t *testing.T) {
	router, _, _ := setupTestImageHandler(t)

	req, _ := http.NewRequest(http.MethodGet, "/images/nonexistent.png", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestImageHandler_DeleteImage_NotFound(t *testing.T) {
	router, _, _ := setupTestImageHandler(t)

	req, _ := http.NewRequest(http.MethodDelete, "/images/nonexistent.png", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
