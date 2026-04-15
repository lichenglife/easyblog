package db

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// RunMigrations 执行数据库迁移脚本
func (db *DB) RunMigrations(sqlFile string) error {
	// 读取 SQL 文件
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取迁移文件失败：%v", err)
	}

	// 执行 SQL 脚本
	if err := db.Exec(string(content)).Error; err != nil {
		return fmt.Errorf("执行迁移脚本失败：%v", err)
	}

	return nil
}

// RunMigrationsFromDir 从目录执行所有 SQL 迁移文件
func (db *DB) RunMigrationsFromDir(dir string) error {
	// 读取目录下的所有 SQL 文件
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("读取迁移目录失败：%v", err)
	}

	for _, file := range files {
		// 读取 SQL 文件
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("读取迁移文件 %s 失败：%v", file, err)
		}

		// 执行 SQL 脚本
		if err := db.DB.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("执行迁移脚本 %s 失败：%v", file, err)
		}
	}

	return nil
}

// ReadFromReader 从 io.Reader 执行 SQL
func (db *DB) ReadFromReader(r io.Reader) error {
	content, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("读取 SQL 内容失败：%v", err)
	}

	if err := db.Exec(string(content)).Error; err != nil {
		return fmt.Errorf("执行 SQL 失败：%v", err)
	}

	return nil
}
