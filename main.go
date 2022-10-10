package main

import (
	"fmt"
	"thub/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	// 注册路由
	bootstrap.SetupRoute(router)

	err := router.Run(":9999")
	if err != nil {
		fmt.Println(err.Error())
	}
}
