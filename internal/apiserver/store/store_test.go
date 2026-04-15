package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDataStore_User(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Act
	store := NewStore(db)

	// Assert
	assert.NotNil(t, store)
	assert.NotNil(t, store.User())
}

func TestDataStore_Post(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Act
	store := NewStore(db)

	// Assert
	assert.NotNil(t, store)
	assert.NotNil(t, store.Post())
}

func TestDataStore_Close(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	store := NewStore(db)

	// Act
	err = store.Close()

	// Assert
	assert.NoError(t, err)
}

func TestNewStore_Singleton(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Act
	store1 := NewStore(db)
	store2 := NewStore(db)

	// Assert
	// 由于使用了 sync.Once，两次调用应该返回同一个实例
	assert.Equal(t, store1, store2)
}
