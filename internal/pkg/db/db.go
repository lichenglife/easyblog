package db

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 表示数据库实例
type DB struct {
	*gorm.DB
}

// NewDB 创建一个新的数据库实例
func NewDB(config *viper.Viper) (*DB, error) {
	// 构建DSN
	username := config.GetString("db.username")
	password := config.GetString("db.password")
	host := config.GetString("db.host")
	port := config.GetInt("db.port")
	database := config.GetString("db.database")

	// 打印调试信息，确认配置值
	fmt.Printf("数据库连接配置: host=%s, port=%d, user=%s, database=%s\n",
		host, port, username, database)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(config.GetString("db.logLevel"))),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库对象失败: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(config.GetInt("db.maxIdleConns"))
	sqlDB.SetMaxOpenConns(config.GetInt("db.maxOpenConns"))
	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("db.connMaxLifetime")) * time.Second)

	return &DB{db}, nil
}

// getLogLevel 获取日志级别
func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库对象失败: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %v", err)
	}

	return nil
}
