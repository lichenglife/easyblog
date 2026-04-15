package store

import (
	"context"
	"testing"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupLikesTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.Like{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestLikes_Create_Success(t *testing.T) {
	// Arrange
	db := setupLikesTestDB(t)
	likeStore := NewLikes(db)

	ctx := context.Background()
	testLike := &model.Like{
		UserID:     1,
		TargetID:   100,
		TargetType: 1, // 文章
	}

	// Act
	err := likeStore.Create(ctx, testLike)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, testLike.ID)

	// Verify
	var savedLike model.Like
	result := db.First(&savedLike, testLike.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, uint(1), savedLike.UserID)
	assert.Equal(t, uint(100), savedLike.TargetID)
	assert.Equal(t, 1, savedLike.TargetType)
}

func TestLikes_Delete_Success(t *testing.T) {
	// Arrange
	db := setupLikesTestDB(t)
	likeStore := NewLikes(db)

	ctx := context.Background()
	testLike := &model.Like{
		UserID:     1,
		TargetID:   100,
		TargetType: 1,
	}
	db.Create(testLike)

	// Act
	err := likeStore.Delete(ctx, testLike.UserID, testLike.TargetID, uint(testLike.TargetType))

	// Assert
	assert.NoError(t, err)

	var deletedLike model.Like
	result := db.First(&deletedLike, testLike.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestLikes_Get_Success(t *testing.T) {
	// Arrange
	db := setupLikesTestDB(t)
	likeStore := NewLikes(db)

	ctx := context.Background()
	testLike := &model.Like{
		UserID:     1,
		TargetID:   100,
		TargetType: 1,
	}
	db.Create(testLike)

	// Act
	result, err := likeStore.Get(ctx, testLike.UserID, testLike.TargetID, uint(testLike.TargetType))

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(100), result.TargetID)
}

func TestLikes_Get_NotFound(t *testing.T) {
	// Arrange
	db := setupLikesTestDB(t)
	likeStore := NewLikes(db)

	ctx := context.Background()

	// Act
	result, err := likeStore.Get(ctx, 999, 999, 1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestLikes_IsLiked_True(t *testing.T) {
	// Arrange
	db := setupLikesTestDB(t)
	likeStore := NewLikes(db)

	ctx := context.Background()
	testLike := &model.Like{
		UserID:     1,
		TargetID:   100,
		TargetType: 1,
	}
	db.Create(testLike)

	// Act
	liked, err := likeStore.IsLiked(ctx, testLike.UserID, testLike.TargetID, uint(testLike.TargetType))

	// Assert
	assert.NoError(t, err)
	assert.True(t, liked)
}

func TestLikes_IsLiked_False(t *testing.T) {
	// Arrange
	db := setupLikesTestDB(t)
	likeStore := NewLikes(db)

	ctx := context.Background()

	// Act
	liked, err := likeStore.IsLiked(ctx, 999, 999, 1)

	// Assert
	assert.NoError(t, err)
	assert.False(t, liked)
}

func TestLikes_CountByTarget_Success(t *testing.T) {
	// Arrange
	db := setupLikesTestDB(t)
	likeStore := NewLikes(db)

	ctx := context.Background()

	// Create multiple likes for the same target
	for i := uint(1); i <= 5; i++ {
		testLike := &model.Like{
			UserID:     i,
			TargetID:   100,
			TargetType: 1,
		}
		db.Create(testLike)
	}

	// Act
	count, err := likeStore.CountByTarget(ctx, 100, 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestLikes_CountByTarget_Zero(t *testing.T) {
	// Arrange
	db := setupLikesTestDB(t)
	likeStore := NewLikes(db)

	ctx := context.Background()

	// Act
	count, err := likeStore.CountByTarget(ctx, 999, 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
