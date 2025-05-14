package main

import (
	"fmt"
	"os"

	"github.com/lichenglife/easyblog/cmd/apiserver/app"
)

func main() {

	command := app.NewAPIServerCommand()

	if err := command.Execute(); err != nil {
		fmt.Printf("执行命令失败: %v\n", err)
		os.Exit(1)
	}

}
