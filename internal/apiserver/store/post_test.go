package store

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lichenglife/easyblog/internal/apiserver/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupPostMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {

	// mock数据库连接db 、用于管理返回SQL执行结构的mock
	// New creates sqlmock database connection and a mock to manage expectations.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	// 初始化gorm
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm %v", err)
	}

	// 关闭数据库资源
	cleanUp := func() {
		defer db.Close()
	}
	return gdb, mock, cleanUp

}
func TestPosts_Create(t *testing.T) {
	db, mock, cleanup := setupPostMockDB(t)
	defer cleanup()
	store := NewPosts(db)

	post := &model.Post{PostID: "p1", UserID: "u1", Title: "test", Content: "content", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO")).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := store.Create(context.Background(), post)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPosts_GetByID(t *testing.T) {
	db, mock, cleanup := setupPostMockDB(t)
	defer cleanup()
	store := NewPosts(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `post` WHERE id = ? ORDER BY `post`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "postID", "userID", "title", "content", "created_at", "updated_at"}).AddRow(1, "p1", "u1", "test", "content", time.Now(), time.Now()))

	post, err := store.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "p1", post.PostID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPosts_Update(t *testing.T) {
	db, mock, cleanup := setupPostMockDB(t)
	defer cleanup()
	store := NewPosts(db)

	post := &model.Post{PostID: "p1", UserID: "u1", Title: "test", Content: "content"}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `post` SET")).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := store.Update(context.Background(), post)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPosts_Delete(t *testing.T) {
	db, mock, cleanup := setupPostMockDB(t)
	defer cleanup()
	store := NewPosts(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `post` WHERE postID = ?")).WithArgs("p1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := store.Delete(context.Background(), "p1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPosts_List(t *testing.T) {
	db, mock, cleanup := setupPostMockDB(t)
	defer cleanup()
	store := NewPosts(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `post`")).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `post` LIMIT ?")).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "postID", "userID", "title", "content", "created_at", "updated_at"}).
			AddRow(1, "p1", "u1", "test", "content", time.Now(), time.Now()))

	count, posts, err := store.List(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Len(t, posts, 1)
	assert.Equal(t, "p1", posts[0].PostID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestPosts_GetByUserID(t *testing.T) {
	db, mock, cleanup := setupPostMockDB(t)
	defer cleanup()
	store := NewPosts(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `post` WHERE userID = ?")).
		WithArgs("u1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `post` WHERE userID = ? LIMIT ? ")).
		WithArgs("u1", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "postID", "userID", "title", "content", "created_at", "updated_at"}).
			AddRow(1, "p1", "u1", "test", "content", time.Now(), time.Now()))

	count, posts, err := store.GetByUserID(context.Background(), "u1", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Len(t, posts, 1)
	assert.Equal(t, "p1", posts[0].PostID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPosts_GetByPostID(t *testing.T) {
	db, mock, cleanup := setupPostMockDB(t)
	defer cleanup()
	store := NewPosts(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `post` WHERE postID = ? ORDER BY `post`.`id` LIMIT ?")).WithArgs("p1", 1).WillReturnRows(sqlmock.NewRows([]string{"id", "postID", "userID", "title", "content", "created_at", "updated_at"}).AddRow(1, "p1", "u1", "test", "content", time.Now(), time.Now()))

	post, err := store.GetByPostID(context.Background(), "p1")
	assert.NoError(t, err)
	assert.Equal(t, "p1", post.PostID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
