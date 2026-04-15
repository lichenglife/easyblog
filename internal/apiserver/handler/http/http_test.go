package handler

import (
	"testing"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDBWithCategories creates an in-memory SQLite database with all tables migrated
// Each call creates a NEW database instance to ensure test isolation
func setupTestDBWithCategories(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Auto migrate all tables
	err = db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
		&model.Like{},
		&model.Category{},
		&model.Tag{},
		&model.Image{},
		&model.SystemConfig{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
