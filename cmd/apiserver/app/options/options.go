package options

import "github.com/spf13/pflag"

// Options 定义全局变量

type Options struct {

	// ServerOptions 定义服务器选项
	ServerOpts *ServerOptions

	// LogOptions 定义日志选项
	LogOpts *LogOptions
	// DBOptions 定义数据库选项
	DBOpts *DBOptions
	// RedisOptions 定义Redis选项
	CacheOpts *CacheOpts
}

// NewOptions 创建默认选项
func NewOptions() *Options {
	return &Options{
		ServerOpts: NewServerOptions(),
		LogOpts:    NewLogOptions(),
		DBOpts:     NewDBOptions(),
		CacheOpts:  NewCacheOptions(),
	}
}

// 添加所有命令行标志
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.ServerOpts.AddFlags(fs)
	o.LogOpts.AddFlags(fs)
	o.DBOpts.AddFlags(fs)
	o.CacheOpts.AddFlags(fs)
}

// Complete 完成选项\
func (o *Options) Complete() error {
	if err := o.ServerOpts.Complete(); err != nil {
		return err
	}
	if err := o.LogOpts.Complete(); err != nil {
		return err
	}
	if err := o.DBOpts.Complete(); err != nil {
		return err
	}
	if err := o.CacheOpts.Complete(); err != nil {
		return err
	}
	return nil
}

// Validate 验证参数
func (o *Options) Validate() error {
	if err := o.ServerOpts.Validate(); err != nil {
		return err
	}
	if err := o.LogOpts.Validate(); err != nil {
		return err
	}
	if err := o.DBOpts.Validate(); err != nil {
		return err
	}
	if err := o.CacheOpts.Validate(); err != nil {
		return err
	}

	return nil

}

// 定义服务参数
type ServerOptions struct {
	// 配置文件路径
	ConfigFile string
	// 服务端口
	Port int
	// 服务模式
	Mode string
	// ReadTimeout 读取超时时间
	ReadTimeout int
	// WriteTimeout 写入超时时间
	WriteTimeout int
	// MaxHeaderBytes 最大请求头大小
	MaxHeaderBytes int
}

// NewServerOptions 创建默认服务选项
func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		ConfigFile:     "../../configs/apiserver.yaml",
		Port:           8080,
		Mode:           "debug",
		ReadTimeout:    60,
		WriteTimeout:   60,
		MaxHeaderBytes: 1 << 20,
	}
}

// Addflags 添加命令行标志
func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.ConfigFile, "config", "c", o.ConfigFile, "config file path")
	fs.IntVarP(&o.Port, "port", "p", o.Port, "server port")
	fs.StringVarP(&o.Mode, "mode", "m", o.Mode, "server mode")
	fs.IntVarP(&o.ReadTimeout, "read-timeout", "r", o.ReadTimeout, "server read timeout")
	fs.IntVarP(&o.WriteTimeout, "write-timeout", "w", o.WriteTimeout, "server write timeout")
	fs.IntVarP(&o.MaxHeaderBytes, "max-header-bytes", "b", o.MaxHeaderBytes, "server max header bytes")
}

// Complete 完成选项
func (o *ServerOptions) Complete() error {
	return nil
}

// Validate 验证选项
func (o *ServerOptions) Validate() error {
	return nil
}

type LogOptions struct {
	// 日志级别
	Level string
	// 日志格式
	Format string
	// 日志文件路径
	Dir string
	// MaxBackups 最大备份数
	MaxBackups int
	// MaxSize  日志文件最大(Mb)
	MaxSize int
	// MaxAge 最大保留时间（天）
	MaxAge int
	// Compress 是否压缩
	Compress bool
}

// NewLogoptions 创建默认日志选项
func NewLogOptions() *LogOptions {
	return &LogOptions{
		Level:      "info",
		Format:     "text",
		Dir:        "logs",
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     3,
		Compress:   false,
	}
}

// AddFlags 添加命令行标志
func (o *LogOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Level, "log-level", "l", o.Level, "日志级别 (debug, info, warn, error, fatal)")
	fs.StringVar(&o.Dir, "log-dir", o.Dir, "日志目录")
	fs.IntVar(&o.MaxSize, "log-max-size", o.MaxSize, "最大日志大小(MB)")
	fs.IntVar(&o.MaxBackups, "log-max-backups", o.MaxBackups, "最大备份数")
	fs.IntVar(&o.MaxAge, "log-max-age", o.MaxAge, "最大保留时间(天)")
	fs.BoolVar(&o.Compress, "log-compress", o.Compress, "是否压缩")
}

// Complete
func (o *LogOptions) Complete() error {
	return nil
}

// Validate 验证选项
func (o *LogOptions) Validate() error {
	return nil
}

type DBOptions struct {
	// Host 数据库主机
	Host string
	// Port  数据库端口
	Port int
	// Username 用户名
	Username string
	// Password 密码
	Password string
	// Database  数据库名
	Database string
	// LogLevel 日志级别
	LogLevel string
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int
	// ConnMaxLifetime 连接最大生命周期(秒)
	ConnMaxLifetime int
}

// NewDBOptions 默认数据库配置

func NewDBOptions() *DBOptions {
	return &DBOptions{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "root",
		Database:        "apiserver",
		LogLevel:        "info",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 3600,
	}
}

// AddFlags 设置命令行标志

func (o *DBOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "db-host", o.Host, "数据库主机")
	fs.IntVar(&o.Port, "db-port", o.Port, "数据库端口")
	fs.StringVar(&o.Username, "db-username", o.Username, "数据库用户名")
	fs.StringVar(&o.Password, "db-password", o.Password, "数据库密码")
	fs.StringVar(&o.Database, "db-database", o.Database, "数据库名")
	fs.StringVar(&o.LogLevel, "db-log-level", o.LogLevel, "数据库日志级别 (silent, error, warn, info)")
	fs.IntVar(&o.MaxIdleConns, "db-max-idle-conns", o.MaxIdleConns, "最大空闲连接数")
	fs.IntVar(&o.MaxOpenConns, "db-max-open-conns", o.MaxOpenConns, "最大打开连接数")
	fs.IntVar(&o.ConnMaxLifetime, "db-conn-max-lifetime", o.ConnMaxLifetime, "连接最大生命周期(秒)")
}

// Complete 完成选项
func (o *DBOptions) Complete() error {
	return nil
}

// Validate 验证选项
func (o *DBOptions) Validate() error {
	return nil
}

type CacheOpts struct {

	// Host 缓存主机
	Host string
	// Port 缓存端口
	Port int
	// Password 缓存密码
	Password string
	// DB 缓存数据库
	DB int
	// PoolSize 连接池大小
	PoolSize int
	// MinIdleConns 最小空闲连接数
	MinIdleConns int
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int
}

// NewCacheoptions  创建默认缓存配置对象
func NewCacheOptions() *CacheOpts {
	return &CacheOpts{
		Host:         "127.0.0.1",
		Port:         6379,
		Password:     "root",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 10,
		MaxIdleConns: 100,
	}
}

// AddFlags 设置默认命令行标志
func (o *CacheOpts) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "cache-host", o.Host, "缓存主机")
	fs.IntVar(&o.Port, "cache-port", o.Port, "缓存端口")
	fs.StringVar(&o.Password, "cache-password", o.Password, "缓存密码")
	fs.IntVar(&o.DB, "cache-db", o.DB, "缓存数据库")
	fs.IntVar(&o.PoolSize, "cache-pool-size", o.PoolSize, "连接池大小")
	fs.IntVar(&o.MinIdleConns, "cache-min-idle-conns", o.MinIdleConns, "最小空闲连接数")
}

// Complete
func (o *CacheOpts) Complete() error {
	return nil
}

// Validate 验证选项
func (o *CacheOpts) Validate() error {
	return nil
}
