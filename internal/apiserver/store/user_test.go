package store

import (
	"context"
	"testing"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestUsers_Create_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()
	testUser := &model.User{
		UserID:   "test-user-001",
		Username: "testuser",
		Password: "hashed_password",
		NickName: "Test User",
		Email:    "test@example.com",
		Phone:    "13800138000",
	}

	// Act
	err := userStore.Create(ctx, testUser)

	// Assert
	assert.NoError(t, err)

	// Verify by reading from database
	var savedUser model.User
	result := db.First(&savedUser, testUser.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, testUser.Username, savedUser.Username)
	assert.Equal(t, testUser.NickName, savedUser.NickName)
	assert.Equal(t, testUser.Email, savedUser.Email)
}

func TestUsers_GetByUsername_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()
	testUser := &model.User{
		UserID:   "test-user-002",
		Username: "finduser",
		Password: "hashed_password",
		NickName: "Find User",
		Email:    "find@example.com",
		Phone:    "13800138001",
	}
	db.Create(testUser)

	// Act
	result, err := userStore.GetByUsername(ctx, "finduser")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "finduser", result.Username)
	assert.Equal(t, "Find User", result.NickName)
}

func TestUsers_GetByUsername_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()

	// Act
	result, err := userStore.GetByUsername(ctx, "nonexistent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUsers_GetByID_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()
	testUser := &model.User{
		UserID:   "test-user-003",
		Username: "iduser",
		Password: "hashed_password",
		NickName: "ID User",
		Email:    "id@example.com",
		Phone:    "13800138002",
	}
	db.Create(testUser)

	// Act
	result, err := userStore.GetByID(ctx, testUser.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "iduser", result.Username)
}

func TestUsers_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()

	// Act
	result, err := userStore.GetByID(ctx, 99999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUsers_Update_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()
	testUser := &model.User{
		UserID:   "test-user-004",
		Username: "updateuser",
		Password: "hashed_password",
		NickName: "Original Name",
		Email:    "original@example.com",
		Phone:    "13800138003",
	}
	db.Create(testUser)

	// Act
	testUser.NickName = "Updated Name"
	testUser.Email = "updated@example.com"
	err := userStore.Update(ctx, testUser)

	// Assert
	assert.NoError(t, err)

	var updatedUser model.User
	result := db.First(&updatedUser, "id = ?", testUser.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Updated Name", updatedUser.NickName)
	assert.Equal(t, "updated@example.com", updatedUser.Email)
}

func TestUsers_Delete_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()
	testUser := &model.User{
		UserID:   "test-user-005",
		Username: "deleteuser",
		Password: "hashed_password",
		NickName: "Delete User",
		Email:    "delete@example.com",
		Phone:    "13800138004",
	}
	db.Create(testUser)

	// Act
	err := userStore.Delete(ctx, testUser.ID)

	// Assert
	assert.NoError(t, err)

	var deletedUser model.User
	result := db.First(&deletedUser, "id = ?", testUser.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestUsers_List_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()

	// Create test data
	for i := 0; i < 15; i++ {
		testUser := &model.User{
			UserID:   "test-user-" + string(rune('a'+i)),
			Username: "user" + string(rune('a'+i)),
			Password: "hashed_password",
			NickName: "User " + string(rune('a'+i)),
			Email:    "user" + string(rune('a'+i)) + "@example.com",
			Phone:    "1380013800" + string(rune('5'+i)),
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		}
		db.Create(testUser)
	}

	// Act
	users, err := userStore.List(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, users, 10)
}

func TestUsers_List_EmptyResult(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()

	// Act
	users, err := userStore.List(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, users)
}

func TestUsers_Create_DuplicateUsername(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	userStore := NewUsers(db)

	ctx := context.Background()
	testUser1 := &model.User{
		UserID:   "test-user-dup1",
		Username: "duplicateuser",
		Password: "hashed_password",
		NickName: "Duplicate User 1",
		Email:    "dup1@example.com",
		Phone:    "13800138100",
	}
	testUser2 := &model.User{
		UserID:   "test-user-dup2",
		Username: "duplicateuser",
		Password: "hashed_password",
		NickName: "Duplicate User 2",
		Email:    "dup2@example.com",
		Phone:    "13800138101",
	}

	// Act
	err1 := userStore.Create(ctx, testUser1)
	err2 := userStore.Create(ctx, testUser2)

	// Assert
	assert.NoError(t, err1)
	assert.Error(t, err2) // Should fail due to unique constraint
}
