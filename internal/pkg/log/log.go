package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 表示日志实例
type Logger struct {
	*zap.Logger
}

// NewLogger 创建一个新的日志实例
func NewLogger(config *viper.Viper) (*Logger, error) {
	// 创建日志目录
	logDir := config.GetString("log.dir")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建日志文件
	logFile := filepath.Join(logDir, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	writer := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    config.GetInt("log.maxSize"),
		MaxBackups: config.GetInt("log.maxBackups"),
		MaxAge:     config.GetInt("log.maxAge"),
		Compress:   config.GetBool("log.compress"),
	}

	// 创建编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writer)),
		zap.NewAtomicLevelAt(getLogLevel(config.GetString("log.level"))),
	)

	// 创建日志实例
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return &Logger{logger}, nil
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug 输出调试日志
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Info 输出信息日志
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Warn 输出警告日志
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Error 输出错误日志
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal 输出致命错误日志
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Sync 同步日志
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// With 创建带有字段的日志对象
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}
