package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPostsTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.Post{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestPosts_Create_Success(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()
	testPost := &model.Post{
		UserID:  1,
		PostID:  "post-001",
		Title:   "Test Post Title",
		Content: "Test post content",
	}

	// Act
	err := postStore.Create(ctx, testPost)

	// Assert
	assert.NoError(t, err)

	// Verify
	var savedPost model.Post
	result := db.First(&savedPost, "post_id = ?", testPost.PostID)
	assert.NoError(t, result.Error)
	assert.Equal(t, testPost.Title, savedPost.Title)
	assert.Equal(t, testPost.Content, savedPost.Content)
}

func TestPosts_GetByID_Success(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()
	testPost := &model.Post{
		UserID:  2,
		PostID:  "post-002",
		Title:   "Find Post",
		Content: "Find content",
	}
	db.Create(testPost)

	// Act
	result, err := postStore.GetByID(ctx, testPost.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Find Post", result.Title)
}

func TestPosts_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()

	// Act
	result, err := postStore.GetByID(ctx, 99999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestPosts_GetByPostID_Success(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()
	testPost := &model.Post{
		UserID:  3,
		PostID:  "post-003",
		Title:   "Get By PostID",
		Content: "Get by postID content",
	}
	db.Create(testPost)

	// Act
	result, err := postStore.GetByPostID(ctx, testPost.PostID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testPost.PostID, result.PostID)
}

func TestPosts_Update_Success(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()
	testPost := &model.Post{
		UserID:  4,
		PostID:  "post-004",
		Title:   "Original Title",
		Content: "Original content",
	}
	db.Create(testPost)

	// Act
	testPost.Title = "Updated Title"
	testPost.Content = "Updated content"
	err := postStore.Update(ctx, testPost)

	// Assert
	assert.NoError(t, err)

	var updatedPost model.Post
	result := db.First(&updatedPost, "id = ?", testPost.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated Title", updatedPost.Title)
}

func TestPosts_Delete_Success(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()
	testPost := &model.Post{
		UserID:  5,
		PostID:  "post-005",
		Title:   "Delete Post",
		Content: "Delete content",
	}
	db.Create(testPost)

	// Act
	err := postStore.Delete(ctx, testPost.ID)

	// Assert
	assert.NoError(t, err)

	var deletedPost model.Post
	result := db.First(&deletedPost, "id = ?", testPost.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestPosts_List_Success(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()

	// Create test data
	for i := uint(1); i <= 15; i++ {
		testPost := &model.Post{
			UserID:  i,
			PostID:  fmt.Sprintf("post-%03d", i),
			Title:   fmt.Sprintf("Post %d", i),
			Content: fmt.Sprintf("Content %d", i),
		}
		db.Create(testPost)
	}

	// Act
	posts, err := postStore.List(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, posts, 10)
}

func TestPosts_List_ByUserID(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()

	// Create test data with same user
	for i := 0; i < 5; i++ {
		testPost := &model.Post{
			UserID:  100,
			PostID:  fmt.Sprintf("post-user-%d", i),
			Title:   fmt.Sprintf("User Post %d", i),
			Content: fmt.Sprintf("User Content %d", i),
		}
		db.Create(testPost)
	}

	// Act
	posts, err := postStore.ListByUserID(ctx, 100, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, posts, 5)
}

func TestPosts_List_EmptyResult(t *testing.T) {
	// Arrange
	db := setupPostsTestDB(t)
	postStore := NewPosts(db)

	ctx := context.Background()

	// Act
	posts, err := postStore.List(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, posts)
}
