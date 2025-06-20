package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
)

func TestNewDBWithDialector(t *testing.T) {
	// 1、设置配置参数
	config := viper.New()
	config.Set("db.host", "localhost")
	config.Set("db.port", "6379")
	config.Set("db.password", "123456")
	config.Set("db.db", "0")
	config.Set("db.logLevel", "info")
	// 2、mocksql 客户端
	sqlDB, _, err := sqlmock.New()
	assert.NoError(t, err)

	// 3、使用mock 创建Gorm  Dialector
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	// 4、调用被测试函数，注入mock的Dialector
	db, err := newDBWithDialector(dialector, config)

	// 5、断言
	assert.NoError(t, err)
	assert.NotNil(t, db)

}

func TestNewDBWithDialector_Failed(t *testing.T) {
	// 1、设置配置参数
	config := viper.New()
	config.Set("db.host", "localhost")
	config.Set("db.port", "6379")
	config.Set("db.password", "123456")
	config.Set("db.db", "0")
	config.Set("db.logLevel", "info")
	db, err := newDBWithDialector(nil, config)

	// 5、断言
	assert.Error(t, err)
	assert.Nil(t, db)

}
