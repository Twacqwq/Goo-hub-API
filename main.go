package main

import (
	"thub/app/cmd"
	initConfig "thub/config"
)

func init() {
	// 加载 config 包的配置信息
	initConfig.Initialize()
}

func main() {
	cmd.Execute()
}
