package main

import (
	"log"
	"os"

	"github.com/lichenglife/easyblog/cmd/apiserver/app"
	_ "github.com/lichenglife/easyblog/docs" // 导入 Swagger 文档
)

// @swagger 2.0
// @title EasyBlog API
// @version 1.0
// @description EasyBlog 是一个简单的博客系统 API 文档。
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey  Bearer
// @in                         header
// @name                       Authorization
// @description               Bearer token authentication
func main() {
	command := app.NewAPIServerCommand()

	if err := command.Execute(); err != nil {
		log.Printf("执行命令失败: %v\n", err)
		os.Exit(1)
	}
}
