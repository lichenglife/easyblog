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

// mock 数据库
func setupUserMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {

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
func TestNewUsers(t *testing.T) {
	//  mock db and  mock manager
	db, mock, cleanUp := setupUserMockDB(t)
	defer cleanUp()
	// 创建store
	store := NewUsers(db)
	// manager  mock
	user := &model.User{
		UserID:    "user-1",
		Username:  "zhangsan",
		Password:  "test",
		NickName:  "zhangsan",
		Email:     "70558555@qq.com",
		Phone:     "1585668545",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO")).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := store.Create(context.Background(), user)
	if err != nil {
		t.Errorf("create failed user %v", err)
	}
	// 确保所有SQL 执行
	mock.ExpectationsWereMet()
	// 验证期待
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func Test_users_GetByID(t *testing.T) {
	// mock gorm   mockmanager
	db, mock, cleanUp := setupUserMockDB(t)
	defer cleanUp()
	// 创建store
	store := NewUsers(db)
	// manager mock
	user := &model.User{
		UserID:    "user-1",
		Username:  "zhangsan",
		Password:  "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		NickName:  "zhangsan",
		Email:     "70558555@qq.com",
		Phone:     "1585668545",
	}

	// 创建模拟的行数据
	rows := sqlmock.NewRows([]string{"userID", "username", "password", "createdAt", "updatedAt", "nickname", "email", "phone"}).
		AddRow(user.UserID, user.Username, user.Password, user.CreatedAt, user.UpdatedAt, user.NickName, user.Email, user.Phone)
	originalQuery := "SELECT * FROM `user` WHERE userID =? ORDER BY `user`.`id` LIMIT ?"
	expectedQuery := regexp.QuoteMeta(originalQuery)
	mock.ExpectQuery(expectedQuery).
		WithArgs(user.UserID, 1).
		WillReturnRows(rows)
	// 执行被测试函数
	got, err := store.GetByID(context.Background(), user.UserID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, user, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_users_GetByID_Failed(t *testing.T) {
	// mock gorm   mockmanager
	db, mock, cleanUp := setupUserMockDB(t)
	defer cleanUp()
	// 创建store
	store := NewUsers(db)
	// manager mock

	// 创建模拟的行数据
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user` WHERE userID =? ORDER BY `user`.`id` LIMIT ?")).
		WithArgs("user-1", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// 执行被测试函数
	_, err := store.GetByID(context.Background(), "user-1")

	// 验证结果
	assert.Error(t, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func Test_users_GetByUsername(t *testing.T) {
	// mock gorm   mockmanager
	db, mock, cleanUp := setupUserMockDB(t)
	defer cleanUp()
	// 创建store
	store := NewUsers(db)
	// manager mock
	user := &model.User{
		UserID:    "user-1",
		Username:  "zhangsan",
		Password:  "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		NickName:  "zhangsan",
		Email:     "70558555@qq.com",
		Phone:     "1585668545",
	}

	// 创建模拟的行数据
	rows := sqlmock.NewRows([]string{"userID", "username", "password", "createdAt", "updatedAt", "nickname", "email", "phone"}).
		AddRow(user.UserID, user.Username, user.Password, user.CreatedAt, user.UpdatedAt, user.NickName, user.Email, user.Phone)
	originalQuery := "SELECT * FROM `user` WHERE username = ? ORDER BY `user`.`id` LIMIT ?"
	expectedQuery := regexp.QuoteMeta(originalQuery)
	t.Logf("Original SQL query: %s", originalQuery)

	t.Logf("Escaped regexp for SQL query: %s", expectedQuery)

	mock.ExpectQuery(expectedQuery).
		WithArgs(user.Username, 1).
		WillReturnRows(rows)
	// 执行被测试函数
	got, err := store.GetByUsername(context.Background(), user.Username)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, user, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_users_Update(t *testing.T) {
	db, mock, cleanUp := setupUserMockDB(t)
	defer cleanUp()

	store := NewUsers(db)

	//  mock store
	// manager mock
	user := &model.User{
		UserID:    "user-1",
		Username:  "zhangsan",
		Password:  "test",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		NickName:  "zhangsan",
		Email:     "70558555@qq.com",
		Phone:     "1585668545",
	}
	// mock store
	// sqlmock.AnyArg() 匹配任意数据格式
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `user` SET `userID`=?,`username`=?,`password`=?,`nickname`=?,`email`=?,`phone`=?,`createdAt`=?,`updatedAt`=? WHERE userID = ?")).
		WithArgs(user.UserID, user.Username, user.Password, user.NickName, user.Email, user.Phone, sqlmock.AnyArg(), sqlmock.AnyArg(), user.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := store.Update(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func Test_users_Delete(t *testing.T) {
	db, mock, cleanUp := setupUserMockDB(t)
	defer cleanUp()
	store := NewUsers(db)
	//store := NewStore(db)
	// mock store
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `user` WHERE userID = ?")).
		WithArgs("userID-1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 执行测试函数
	err := store.Delete(context.Background(), "userID-1")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func Test_users_List(t *testing.T) {
	db, mock, cleanUp := setupUserMockDB(t)
	defer cleanUp()
	store := NewUsers(db)
	//store := NewStore(db)

	// mock manager

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `user`")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user`")).
		WillReturnRows(sqlmock.NewRows([]string{"userID", "username", "password", "createdAt", "updatedAt", "nickname", "email", "phone"}).
			AddRow("userID-1", "zhangsan", "test", time.Now(), time.Now(), "zhangsan", "70558555@qq.com", "1585668545"))
	_, _, err := store.List(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
