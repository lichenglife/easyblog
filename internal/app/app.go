package app

import (
	"fmt"

	"github.com/lichenglife/easyblog/internal/apiserver/store"
	"github.com/lichenglife/easyblog/internal/pkg/cache"
	"github.com/lichenglife/easyblog/internal/pkg/db"
	"github.com/lichenglife/easyblog/internal/pkg/log"

	"github.com/spf13/viper"
)

// App 初始化选项
type AppOptions struct {
	EnableDB     bool
	EnableCache  bool
	EnableServer bool
}

// 构造方法， 返回默认配置
func DefaultAppOptions() *AppOptions {

	return &AppOptions{
		EnableDB:     true,
		EnableCache:  false,
		EnableServer: true,
	}
}

// IAPP 定义APP 核心接口， 应用启动初始化的核心方法
type IApp interface {
	GetLogger() *log.Logger
	GetDB() *db.DB
	GetCache() *cache.Cache
	// GetStoreFactory 获取存储工厂实例
	GetStoreFactory() store.Factory

	// 服务接口

	// 关闭应用
	Close() error
}

// App 表示应用实例

type App struct {
	// 配置
	config *viper.Viper
	//  日志
	logger *log.Logger
	//  数据库
	Db *db.DB
	//  缓存
	cache *cache.Cache

	// 服务

	//  存储
	factory store.Factory
	//  认证服务

}

// 创建App实例
func NewAPP(config *viper.Viper) (IApp, error) {

	return NewAppWithOptions(config, DefaultAppOptions())
}

// 使用默认参数构建App实例
func NewAppWithOptions(config *viper.Viper, option *AppOptions) (IApp, error) {

	// 实例化App对象
	app := &App{
		config: config,
	}
	err := app.initLogger()
	if err != nil {
		return nil, fmt.Errorf("初始化日志失败%v", err)
	}
	if option.EnableDB {
		err = app.initDB()
		if err != nil {
			return nil, fmt.Errorf("初始化数据库失败%v", err)
		}
	}
	if option.EnableCache {
		err = app.initCache()
		if err != nil {
			return nil, fmt.Errorf("初始化缓存失败%v", err)
		}
	}
	err = app.initStoreFactory()
	if err != nil {
		return nil, fmt.Errorf("初始化存储工厂失败%v", err)
	}

	return app, nil
}

// initLogger
func (a *App) initLogger() error {

	logger, err := log.NewLogger(a.config)
	if err != nil {
		return err
	}
	a.logger = logger
	return nil
}

// initDB 初始化数据库
func (app *App) initDB() error {
	db, err := db.NewDB(app.config)
	if err != nil {
		return err
	}
	app.Db = db
	return nil
}

func (app *App) initCache() error {
	cache, err := cache.NewCache(app.config)
	if err != nil {
		return err
	}
	app.cache = cache
	return nil
}

func (app *App) initStoreFactory() error {
	app.factory = store.NewFactory(app.Db.DB)
	return nil
}
func (app *App) Close() error {

	if app.Db != nil {
		err := app.Db.Close()
		if err != nil {
			return err
		}

	}
	if app.cache != nil {
		err := app.cache.Close()
		if err != nil {
			return err
		}

	}
	return nil
}

func (app *App) GetLogger() *log.Logger {

	return app.logger
}
func (app *App) GetDB() *db.DB {
	return app.Db
}
func (app *App) GetCache() *cache.Cache {

	return app.cache
}

func (app *App) GetStoreFactory() store.Factory {

	return app.factory
}

// 服务接口
