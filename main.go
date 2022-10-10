package main

import (
	"flag"
	"fmt"
	"thub/bootstrap"
	initConfig "thub/config"
	"thub/pkg/config"

	"github.com/gin-gonic/gin"
)

func init() {
	// 加载 config 包的配置信息
	initConfig.Initialize()
}

func main() {
	// 加载配置
	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件, 如 --env=testing 加载的是 .env.testing 文件")
	flag.Parse()
	config.InitConfig(env)

	// 启动服务器
	router := gin.New()
	bootstrap.SetupRoute(router)
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		fmt.Println(err.Error())
	}
}
