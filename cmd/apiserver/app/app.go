package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lichenglife/easyblog/internal/app"
	"github.com/lichenglife/easyblog/internal/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// AppConfig 应用配置,用于应用初始化

type AppConfig struct {
	// 应用选项
	AppOpts *app.AppOptions
	// 是否启动独立服务器
	EnableServer bool

	// 优雅终止超时时间
	ShutDownTimeout time.Duration
}

func DefaultAppOptions() *AppConfig {

	return &AppConfig{
		AppOpts:         app.DefaultAppOptions(),
		EnableServer:    true,
		ShutDownTimeout: 20 * time.Second,
	}
}

// RunApp 根据配置启动程序
func RunApp(config *viper.Viper) error {

	return RunAppWithDefaultAppOptions(config, DefaultAppOptions())
}

// RunAppWithDefaultAppOptions  启动App应用 启动server服务

func RunAppWithDefaultAppOptions(config *viper.Viper, appConfig *AppConfig) error {

	// 构建app
	app, err := app.NewAppWithOptions(config, appConfig.AppOpts)
	if err != nil {
		return fmt.Errorf("初始化应用失败:%v", err)
	}

	//  构建server
	server, err := server.NewUnionServer(config, app)
	if err != nil {
		return err
	}

	// 初始化server
	if err := server.Init(); err != nil {
		return fmt.Errorf("初始化server失败:%v", err)
	}

	go func() {
		//  启动server
		if err := server.Start(); err != nil {
			app.GetLogger().Error("启动服务失败", zap.Error(err))
			os.Exit(1)
		}
	}()
	app.GetLogger().Info("启动服务成功",
		zap.String("HTTP模式", config.GetString("http.server.mode")),
		zap.Int("HTTP端口", config.GetInt("server.http.port")))
	// 创建一个 os.Signal 类型的 channel，用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号，但 SIGKILL 信号无法被捕获，因此无需监听和处理
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞程序，等待从 quit channel 中接收到信号
	<-quit

	//log.Infow("Shutting down server ...")

	// 优雅关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务，再关闭被依赖的服务
	server.Stop(ctx)

	return nil

}
