package main

import (
	"fmt"
	"os"

	"github.com/lichenglife/easyblog/cmd/apiserver/app"
)
// @title           EasyBlog API
// @version         1.0
// @description     A simple blog server written in Go.
// @termsOfService  http://swagger.io/terms/

// @contact.name   lichenglife
// @contact.url    https://github.com/lichenglife/easyblog
// @contact.email  lichenglife@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  Bearer
// @in                         header
// @name                       Authorization
// @description               Bearer token authentication
func main() {

	command := app.NewAPIServerCommand()

	if err := command.Execute(); err != nil {
		fmt.Printf("执行命令失败: %v\n", err)
		os.Exit(1)
	}

}
