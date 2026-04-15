package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lichenglife/easyblog/internal/apiserver/biz"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/auth"
	"github.com/lichenglife/easyblog/internal/pkg/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestPostHandler(t *testing.T) (*gin.Engine, PostHandler, *viper.Viper) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add mock auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("userID", "1")
		c.Set("userRole", 2) // Admin role
		c.Next()
	})

	config := setupTestConfig()
	logger, _ := log.NewLogger(config)
	db := setupTestDB(t)
	// 使用 NewTestStore 创建独立的 store 实例，避免 once 单例模式导致测试间共享数据
	dbStore := store.NewTestStore(db)
	jwt := auth.NewJWT(config)
	testBiz := biz.NewBiz(dbStore, jwt, logger.Named("test"))
	handler := NewPostHandler(logger, testBiz.PostV1())

	// Setup routes
	v1 := router.Group("/v1")
	{
		v1.POST("/post", handler.CreatePost)
		v1.GET("/post/:id", handler.GetPostByID)
		v1.PUT("/post/:id", handler.UpdatePost)
		v1.DELETE("/post/:id", handler.DeletePost)
		v1.GET("/user/:userID/posts", handler.GetPostsByUserID)
		v1.GET("/posts", handler.ListPosts)
		v1.GET("/posts/search", handler.SearchPosts)
		v1.GET("/user/posts/me", handler.GetMyPosts)
		v1.POST("/admin/post/:id/top", handler.TopPost)
		v1.POST("/admin/post/:id/untop", handler.UnTopPost)
	}

	return router, handler, config
}

// ==================== 创建文章测试 ====================

func TestPostHandler_CreatePost_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	createPostReq := model.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
		Summary: "Test post summary",
	}
	body, _ := json.Marshal(createPostReq)

	// Act
	req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_CreatePost_InvalidParams(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	tests := []struct {
		name       string
		req        model.CreatePostRequest
		wantStatus int
	}{
		{
			name: "标题为空",
			req: model.CreatePostRequest{
				Content: "This is a test post content",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "内容为空",
			req: model.CreatePostRequest{
				Title: "Test Post",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "标题超长",
			req: model.CreatePostRequest{
				Title:   string(make([]byte, 256)),
				Content: "Content",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "摘要超长",
			req: model.CreatePostRequest{
				Title:   "Test Post",
				Content: "Content",
				Summary: string(make([]byte, 501)),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "无效的状态值",
			req: model.CreatePostRequest{
				Title:   "Test Post",
				Content: "Content",
				Status:  5,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.req)
			req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// ==================== 获取文章详情测试 ====================

func TestPostHandler_GetPostByID_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create a post first
	createPostReq := model.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
		Summary: "Test post summary",
	}
	body, _ := json.Marshal(createPostReq)
	req1, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req1)

	// Act - Get post by ID
	req, _ := http.NewRequest(http.MethodGet, "/v1/post/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_GetPostByID_NotFound(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/post/99999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_GetPostByID_InvalidID(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/post/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== 更新文章测试 ====================

func TestPostHandler_UpdatePost_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create a post first
	createPostReq := model.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
		Summary: "Test post summary",
	}
	body, _ := json.Marshal(createPostReq)
	req1, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req1)

	// Update post
	updatePostReq := model.UpdatePostRequest{
		ID:      1,
		Title:   "Updated Post Title",
		Content: "Updated content",
		Summary: "Updated summary",
	}
	updateBody, _ := json.Marshal(updatePostReq)

	// Act
	req, _ := http.NewRequest(http.MethodPut, "/v1/post/1", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_UpdatePost_NotFound(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	updatePostReq := model.UpdatePostRequest{
		ID:      99999,
		Title:   "Updated Title",
		Content: "Updated content",
	}
	body, _ := json.Marshal(updatePostReq)

	// Act
	req, _ := http.NewRequest(http.MethodPut, "/v1/post/99999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_UpdatePost_InvalidParams(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	tests := []struct {
		name       string
		req        model.UpdatePostRequest
		wantStatus int
	}{
		{
			name: "标题为空",
			req: model.UpdatePostRequest{
				ID:      1,
				Content: "Updated content",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "内容为空",
			req: model.UpdatePostRequest{
				ID:    1,
				Title: "Updated Title",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.req)
			req, _ := http.NewRequest(http.MethodPut, "/v1/post/1", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// ==================== 删除文章测试 ====================

func TestPostHandler_DeletePost_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create a post first
	createPostReq := model.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
		Summary: "Test post summary",
	}
	body, _ := json.Marshal(createPostReq)
	req1, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req1)

	// Act - Delete post
	req, _ := http.NewRequest(http.MethodDelete, "/v1/post/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_DeletePost_NotFound(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodDelete, "/v1/post/99999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_DeletePost_InvalidID(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodDelete, "/v1/post/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== 获取文章列表测试 ====================

func TestPostHandler_ListPosts_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create some posts
	for i := 0; i < 5; i++ {
		createPostReq := model.CreatePostRequest{
			Title:   fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("Content %d", i),
			Summary: fmt.Sprintf("Summary %d", i),
		}
		body, _ := json.Marshal(createPostReq)
		req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/posts?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_ListPosts_WithPagination(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create some posts
	for i := 0; i < 25; i++ {
		createPostReq := model.CreatePostRequest{
			Title:   fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("Content %d", i),
			Summary: fmt.Sprintf("Summary %d", i),
		}
		body, _ := json.Marshal(createPostReq)
		req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}

	// Act - Get second page
	req, _ := http.NewRequest(http.MethodGet, "/v1/posts?page=2&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_ListPosts_WithCategoryFilter(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create a post with category
	createPostReq := model.CreatePostRequest{
		Title:      "Categorized Post",
		Content:    "Content with category",
		CategoryID: 1,
	}
	body, _ := json.Marshal(createPostReq)
	req1, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req1)

	// Act - Filter by category
	req, _ := http.NewRequest(http.MethodGet, "/v1/posts?categoryId=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_ListPosts_InvalidCategoryID(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/posts?categoryId=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== 获取用户文章列表测试 ====================

func TestPostHandler_GetPostsByUserID_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/user/1/posts?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_GetPostsByUserID_WithPagination(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create some posts
	for i := 0; i < 15; i++ {
		createPostReq := model.CreatePostRequest{
			Title:   fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("Content %d", i),
		}
		body, _ := json.Marshal(createPostReq)
		req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}

	// Act - Get user posts with pagination
	req, _ := http.NewRequest(http.MethodGet, "/v1/user/1/posts?page=1&pageSize=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_GetPostsByUserID_InvalidUserID(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/user/invalid/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== 搜索文章测试 ====================

func TestPostHandler_SearchPosts_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create a post first
	createPostReq := model.CreatePostRequest{
		Title:   "Test Post for Search",
		Content: "This is a test post content for searching",
		Summary: "Test post summary",
	}
	body, _ := json.Marshal(createPostReq)
	req1, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req1)

	// Act - Search by keyword
	req, _ := http.NewRequest(http.MethodGet, "/v1/posts/search?keyword=Test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_SearchPosts_EmptyKeyword(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/posts/search?keyword=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_SearchPosts_MissingKeyword(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/posts/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostHandler_SearchPosts_WithPagination(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create some posts
	for i := 0; i < 15; i++ {
		createPostReq := model.CreatePostRequest{
			Title:   fmt.Sprintf("Search Test Post %d", i),
			Content: fmt.Sprintf("Content %d", i),
		}
		body, _ := json.Marshal(createPostReq)
		req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}

	// Act - Search with pagination
	req, _ := http.NewRequest(http.MethodGet, "/v1/posts/search?keyword=Search&pageSize=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==================== 我的文章测试 ====================

func TestPostHandler_GetMyPosts_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create a post first
	createPostReq := model.CreatePostRequest{
		Title:   "My Test Post",
		Content: "This is my test post content",
	}
	body, _ := json.Marshal(createPostReq)
	req1, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req1)

	// Act
	req, _ := http.NewRequest(http.MethodGet, "/v1/user/posts/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_GetMyPosts_WithPagination(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create some posts
	for i := 0; i < 15; i++ {
		createPostReq := model.CreatePostRequest{
			Title:   fmt.Sprintf("My Post %d", i),
			Content: fmt.Sprintf("My Content %d", i),
		}
		body, _ := json.Marshal(createPostReq)
		req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}

	// Act - Get my posts with pagination
	req, _ := http.NewRequest(http.MethodGet, "/v1/user/posts/me?page=1&pageSize=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==================== 置顶文章测试 ====================

func TestPostHandler_TopPost_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create a post first
	createPostReq := model.CreatePostRequest{
		Title:   "Post to Top",
		Content: "This post will be topped",
	}
	body, _ := json.Marshal(createPostReq)
	req1, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// Assert create success
	assert.Equal(t, http.StatusOK, w1.Code)

	// Act - Top the post (use ID 1 since it's a fresh database)
	req, _ := http.NewRequest(http.MethodPost, "/v1/admin/post/1/top", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_TopPost_NotFound(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodPost, "/v1/admin/post/99999/top", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_TopPost_InvalidID(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodPost, "/v1/admin/post/invalid/top", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== 取消置顶文章测试 ====================

func TestPostHandler_UnTopPost_Success(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Create a post first
	createPostReq := model.CreatePostRequest{
		Title:   "Post to UnTop",
		Content: "This post will be untopped",
		IsTop:   1,
	}
	body, _ := json.Marshal(createPostReq)
	req1, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// Assert create success
	assert.Equal(t, http.StatusOK, w1.Code)

	// Act - UnTop the post (use ID 1 since it's a fresh database)
	req, _ := http.NewRequest(http.MethodPost, "/v1/admin/post/1/untop", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostHandler_UnTopPost_NotFound(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodPost, "/v1/admin/post/99999/untop", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostHandler_UnTopPost_InvalidID(t *testing.T) {
	router, _, _ := setupTestPostHandler(t)

	// Act
	req, _ := http.NewRequest(http.MethodPost, "/v1/admin/post/invalid/untop", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== 表驱动测试 ====================

func TestPostHandler_CreatePost(t *testing.T) {
	tests := []struct {
		name       string
		req        model.CreatePostRequest
		wantStatus int
	}{
		{
			name: "成功创建文章",
			req: model.CreatePostRequest{
				Title:   "New Post",
				Content: "New content",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "标题为空",
			req: model.CreatePostRequest{
				Content: "Content only",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "内容为空",
			req: model.CreatePostRequest{
				Title: "Title only",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "状态值有效",
			req: model.CreatePostRequest{
				Title:   "Valid Status Post",
				Content: "Content",
				Status:  1,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "置顶值有效",
			req: model.CreatePostRequest{
				Title:   "Topped Post",
				Content: "Content",
				IsTop:   1,
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, _, _ := setupTestPostHandler(t)
			body, _ := json.Marshal(tt.req)
			req, _ := http.NewRequest(http.MethodPost, "/v1/post", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
