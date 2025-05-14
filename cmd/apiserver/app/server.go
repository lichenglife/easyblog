package app

import (
	"fmt"

	"github.com/lichenglife/easyblog/cmd/apiserver/app/options"
	"github.com/spf13/cobra"
)

// NewAPIServerCommand 创建服务启动命令实例

func NewAPIServerCommand() *cobra.Command {
	// 创建命令对象
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "easyblog-apiserver",
		Short: "启动 easyblog服务",
		Long:  `启动 easyblog服务，提供博客管理系统`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 加载配置
			return Run(opts)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			// 验证参数
			for _, arg := range args {
				if len(arg) == 0 {
					return fmt.Errorf("%q 不允许为空", arg)
				}
			}
			return nil
		},
	}

	// 添加命令行参数
	opts.AddFlags(cmd.Flags())

	// 添加子命令
	cmd.AddCommand(NewVersionCommand())

	return cmd
}

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "打印easyblog版本信息",
		Long:  `打印easyblog版本信息`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version: 1.0.0")
		},
	}
}

// Run 启动服务
func Run(opts *options.Options) error {

	// 完成命令行参数加载
	if err := opts.Complete(); err != nil {
		return err
	}

	// 完成命令行参数验证
	if err := opts.Validate(); err != nil {
		return err
	}
	//  加载配置
	cfg, err := LoadConfig(opts)
	if err != nil {
		return fmt.Errorf("初始化配置失败%v", err)
	}

	return RunApp(cfg.Viper)
}
