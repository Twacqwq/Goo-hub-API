package bootstrap

import (
	"net/http"
	"strings"
	"thub/app/http/middlewares"
	"thub/routes"

	"github.com/gin-gonic/gin"
)

// 路由注册中心
func SetupRoute(router *gin.Engine) {

	// 注册全局中间件
	registerGlobalMiddlewarte(router)

	// 注册404请求handler
	setupNotFoundHandler(router)

	// 注册API路由
	routes.RegisterAPIRoutes(router)
}

func registerGlobalMiddlewarte(router *gin.Engine) {
	router.Use(
		middlewares.Logger(),
		middlewares.Recovery(),
	)
}

func setupNotFoundHandler(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		acceptString := c.Request.Header.Get("accept")
		if strings.Contains(acceptString, "text/html") {
			c.String(http.StatusNotFound, "404 NOT FOUND")
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "404 NOT FOUND",
			})
		}
	})
}
