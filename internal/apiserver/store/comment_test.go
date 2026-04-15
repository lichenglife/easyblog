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

func setupCommentsTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.Comment{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestComments_Create_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()
	testComment := &model.Comment{
		PostID:   1,
		UserID:   1,
		Content:  "Test comment content",
		ParentID: 0,
	}

	// Act
	err := commentStore.Create(ctx, testComment)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, testComment.ID)

	// Verify
	var savedComment model.Comment
	result := db.First(&savedComment, testComment.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Test comment content", savedComment.Content)
}

func TestComments_GetByID_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()
	testComment := &model.Comment{
		PostID:   1,
		UserID:   1,
		Content:  "Test comment",
		ParentID: 0,
	}
	db.Create(testComment)

	// Act
	result, err := commentStore.GetByID(ctx, testComment.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test comment", result.Content)
}

func TestComments_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()

	// Act
	result, err := commentStore.GetByID(ctx, 99999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestComments_Update_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()
	testComment := &model.Comment{
		PostID:   1,
		UserID:   1,
		Content:  "Original content",
		ParentID: 0,
	}
	db.Create(testComment)

	// Act
	testComment.Content = "Updated content"
	testComment.IsEdited = 1
	err := commentStore.Update(ctx, testComment)

	// Assert
	assert.NoError(t, err)

	var updatedComment model.Comment
	result := db.First(&updatedComment, testComment.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated content", updatedComment.Content)
	assert.Equal(t, 1, updatedComment.IsEdited)
}

func TestComments_Delete_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()
	testComment := &model.Comment{
		PostID:   1,
		UserID:   1,
		Content:  "Delete me",
		ParentID: 0,
	}
	db.Create(testComment)

	// Act
	err := commentStore.Delete(ctx, testComment.ID)

	// Assert
	assert.NoError(t, err)

	var deletedComment model.Comment
	result := db.First(&deletedComment, testComment.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestComments_ListByPostID_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()

	// Create test data - top level comments
	for i := 0; i < 15; i++ {
		testComment := &model.Comment{
			PostID:   1,
			UserID:   uint(i + 1),
			Content:  fmt.Sprintf("Comment %d", i),
			ParentID: 0,
			Status:   1,
		}
		db.Create(testComment)
	}

	// Act
	comments, err := commentStore.ListByPostID(ctx, 1, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, comments, 10)
}

func TestComments_ListByParentID_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()

	// Create parent comment
	parentComment := &model.Comment{
		PostID:   1,
		UserID:   1,
		Content:  "Parent comment",
		ParentID: 0,
		Status:   1,
	}
	db.Create(parentComment)

	// Create replies
	for i := 0; i < 5; i++ {
		testComment := &model.Comment{
			PostID:   1,
			UserID:   uint(i + 1),
			Content:  fmt.Sprintf("Reply %d", i),
			ParentID: parentComment.ID,
			Status:   1,
		}
		db.Create(testComment)
	}

	// Act
	replies, err := commentStore.ListByParentID(ctx, parentComment.ID, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, replies, 5)
}

func TestComments_IncrementLikeCount_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()
	testComment := &model.Comment{
		PostID:    1,
		UserID:    1,
		Content:   "Test comment",
		LikeCount: 0,
	}
	db.Create(testComment)

	// Act
	err := commentStore.IncrementLikeCount(ctx, testComment.ID)

	// Assert
	assert.NoError(t, err)

	var updatedComment model.Comment
	result := db.First(&updatedComment, testComment.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, uint(1), updatedComment.LikeCount)
}

func TestComments_IncrementReplyCount_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()
	testComment := &model.Comment{
		PostID:     1,
		UserID:     1,
		Content:    "Test comment",
		ReplyCount: 0,
	}
	db.Create(testComment)

	// Act
	err := commentStore.IncrementReplyCount(ctx, testComment.ID)

	// Assert
	assert.NoError(t, err)

	var updatedComment model.Comment
	result := db.First(&updatedComment, testComment.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, uint(1), updatedComment.ReplyCount)
}

func TestComments_DecrementReplyCount_Success(t *testing.T) {
	// Arrange
	db := setupCommentsTestDB(t)
	commentStore := NewComments(db)

	ctx := context.Background()
	testComment := &model.Comment{
		PostID:     1,
		UserID:     1,
		Content:    "Test comment",
		ReplyCount: 3,
	}
	db.Create(testComment)

	// Act
	err := commentStore.DecrementReplyCount(ctx, testComment.ID)

	// Assert
	assert.NoError(t, err)

	var updatedComment model.Comment
	result := db.First(&updatedComment, testComment.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, uint(2), updatedComment.ReplyCount)
}
