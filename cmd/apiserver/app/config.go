package app

import (
	"fmt"

	"github.com/lichenglife/easyblog/cmd/apiserver/app/options"
	"github.com/spf13/viper"
)

// Config is the configuration for the application.
type Config struct {
	ConfigFile string
	Port       int
	Mode       string
	LogLevel   string

	// 配置信息
	*viper.Viper
}

// NewConfig 创建默认配置对象
func NewConfig() *Config {
	return &Config{
		ConfigFile: "configs/apiserver.yaml",
		Port:       8080,
		Mode:       "debug",
		LogLevel:   "info",
	}
}

// LoadConfig 加载配置文件
// 1、支持命令行选项加载配置文件参数
// 2、支持环境变量加载配置文件参数

func LoadConfig(in interface{}) (*Config, error) {

	v := viper.New()

	switch input := in.(type) {
	// 1、支持命令行选项加载配置文件参数
	case *options.Options:
		//  获取命令行选项
		return loadFromCommandLineOptions(v, input)
	case *Config:
		// 从配置文件加载配置
		return loadFromConfigFile(v, input)

	default:
		return nil, fmt.Errorf("unsupported input type: %T", in)

	}

}

// 从配置文件加载配置
func loadFromConfigFile(v *viper.Viper, cfg *Config) (*Config, error) {
	// 设置配置文件地址
	v.SetConfigFile(cfg.ConfigFile)
	// 设置配置文件类型
	v.SetConfigType("yaml")

	// 设置默认值
	v.SetDefault("port", cfg.Port)
	v.SetDefault("mode", cfg.Mode)
	v.SetDefault("log.level", cfg.LogLevel)
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	return &Config{
		ConfigFile: cfg.ConfigFile,
		Port:       v.GetInt("port"),
		Mode:       v.GetString("mode"),
		LogLevel:   v.GetString("log.level"),
		Viper:      v,
	}, nil
}

// 从命令行选项加载配置文件参数
func loadFromCommandLineOptions(v *viper.Viper, opts *options.Options) (*Config, error) {
	v.SetConfigFile(opts.ServerOpts.ConfigFile)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	// 设置默认值
	setDefaultsFromOptions(v, opts)
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// 将命令行参数与配置文件参数合并
	mergeCommandLineOptions(v, opts)
	return &Config{
		ConfigFile: opts.ServerOpts.ConfigFile,
		Port:       v.GetInt("server.http.port"),
		Mode:       v.GetString("server.http.model"),
		// todo 这两种有什么区别？
		LogLevel: v.GetString("log.level"),
		Viper:    v,
	}, nil
}

// 将命令行参数合并到配置中
func mergeCommandLineOptions(v *viper.Viper, opts *options.Options) {
	// HTTP服务器设置
	v.Set("server.http.port", opts.ServerOpts.Port)
	v.Set("server.http.mode", opts.ServerOpts.Mode)
	v.Set("server.http.readTimeout", opts.ServerOpts.ReadTimeout)
	v.Set("server.http.writeTimeout", opts.ServerOpts.WriteTimeout)
	v.Set("server.http.maxHeaderBytes", opts.ServerOpts.MaxHeaderBytes)

	// 日志设置
	v.Set("log.level", opts.LogOpts.Level)
	v.Set("log.dir", opts.LogOpts.Dir)
	v.Set("log.maxSize", opts.LogOpts.MaxSize)
	v.Set("log.maxBackups", opts.LogOpts.MaxBackups)
	v.Set("log.maxAge", opts.LogOpts.MaxAge)
	v.Set("log.compress", opts.LogOpts.Compress)

	// 数据库设置 - 只有当值与默认值不同时，才可能是通过命令行设置的
	defaultOpts := options.NewDBOptions()
	if opts.DBOpts.Username != defaultOpts.Username {
		v.Set("db.username", opts.DBOpts.Username)
	}
	if opts.DBOpts.Password != defaultOpts.Password {
		v.Set("db.password", opts.DBOpts.Password)
	}
	if opts.DBOpts.Host != defaultOpts.Host {
		v.Set("db.host", opts.DBOpts.Host)
	}
	if opts.DBOpts.Port != defaultOpts.Port {
		v.Set("db.port", opts.DBOpts.Port)
	}
	if opts.DBOpts.Database != defaultOpts.Database {
		v.Set("db.database", opts.DBOpts.Database)
	}
	if opts.DBOpts.LogLevel != defaultOpts.LogLevel {
		v.Set("db.logLevel", opts.DBOpts.LogLevel)
	}
	if opts.DBOpts.MaxIdleConns != defaultOpts.MaxIdleConns {
		v.Set("db.maxIdleConns", opts.DBOpts.MaxIdleConns)
	}
	if opts.DBOpts.MaxOpenConns != defaultOpts.MaxOpenConns {
		v.Set("db.maxOpenConns", opts.DBOpts.MaxOpenConns)
	}
	if opts.DBOpts.ConnMaxLifetime != defaultOpts.ConnMaxLifetime {
		v.Set("db.connMaxLifetime", opts.DBOpts.ConnMaxLifetime)
	}

	// 缓存设置 - 同理应用相同的逻辑
	defaultCacheOpts := options.NewCacheOptions()
	if opts.CacheOpts.Host != defaultCacheOpts.Host {
		v.Set("redis.host", opts.CacheOpts.Host)
	}
	if opts.CacheOpts.Port != defaultCacheOpts.Port {
		v.Set("redis.port", opts.CacheOpts.Port)
	}
	if opts.CacheOpts.Password != defaultCacheOpts.Password {
		v.Set("redis.password", opts.CacheOpts.Password)
	}
	if opts.CacheOpts.DB != defaultCacheOpts.DB {
		v.Set("redis.db", opts.CacheOpts.DB)
	}
	if opts.CacheOpts.PoolSize != defaultCacheOpts.PoolSize {
		v.Set("redis.poolSize", opts.CacheOpts.PoolSize)
	}
	if opts.CacheOpts.MinIdleConns != defaultCacheOpts.MinIdleConns {
		v.Set("redis.minIdleConns", opts.CacheOpts.MinIdleConns)
	}
	if opts.CacheOpts.MaxIdleConns != defaultCacheOpts.MaxIdleConns {
		v.Set("redis.maxIdleConns", opts.CacheOpts.MaxIdleConns)
	}
}

func setDefaultsFromOptions(v *viper.Viper, opts *options.Options) {
	// HTTP服务器默认值
	v.SetDefault("server.http.port", opts.ServerOpts.Port)
	v.SetDefault("server.http.mode", opts.ServerOpts.Mode)
	v.SetDefault("server.http.readTimeout", opts.ServerOpts.ReadTimeout)
	v.SetDefault("server.http.writeTimeout", opts.ServerOpts.WriteTimeout)
	v.SetDefault("server.http.maxHeaderBytes", opts.ServerOpts.MaxHeaderBytes)

	// gRPC服务器默认值
	v.SetDefault("server.grpc.port", 9090)
	v.SetDefault("server.grpc.maxConnectionAge", 3600)
	v.SetDefault("server.grpc.maxConnectionAgeGrace", 10)
	v.SetDefault("server.grpc.maxConnectionIdle", 300)
	v.SetDefault("server.grpc.keepAliveTime", 60)
	v.SetDefault("server.grpc.keepAliveTimeout", 20)
	v.SetDefault("server.grpc.maxRecvMsgSize", 4*1024*1024) // 4MB
	v.SetDefault("server.grpc.maxSendMsgSize", 4*1024*1024) // 4MB

	// 日志默认值
	v.SetDefault("log.level", opts.LogOpts.Level)
	v.SetDefault("log.dir", opts.LogOpts.Dir)
	v.SetDefault("log.maxSize", opts.LogOpts.MaxSize)
	v.SetDefault("log.maxBackups", opts.LogOpts.MaxBackups)
	v.SetDefault("log.maxAge", opts.LogOpts.MaxAge)
	v.SetDefault("log.compress", opts.LogOpts.Compress)

	// 数据库默认值
	v.SetDefault("db.host", opts.DBOpts.Host)
	v.SetDefault("db.port", opts.DBOpts.Port)
	v.SetDefault("db.username", opts.DBOpts.Username)
	v.SetDefault("db.password", opts.DBOpts.Password)
	v.SetDefault("db.database", opts.DBOpts.Database)
	v.SetDefault("db.logLevel", opts.DBOpts.LogLevel)
	v.SetDefault("db.maxIdleConns", opts.DBOpts.MaxIdleConns)
	v.SetDefault("db.maxOpenConns", opts.DBOpts.MaxOpenConns)
	v.SetDefault("db.connMaxLifetime", opts.DBOpts.ConnMaxLifetime)

	// 缓存默认值
	v.SetDefault("redis.host", opts.CacheOpts.Host)
	v.SetDefault("redis.port", opts.CacheOpts.Port)
	v.SetDefault("redis.password", opts.CacheOpts.Password)
	v.SetDefault("redis.db", opts.CacheOpts.DB)
	v.SetDefault("redis.poolSize", opts.CacheOpts.PoolSize)
	v.SetDefault("redis.minIdleConns", opts.CacheOpts.MinIdleConns)
	v.SetDefault("redis.maxIdleConns", opts.CacheOpts.MaxIdleConns)
}
