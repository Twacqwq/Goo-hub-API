// Package limiter 处理限流逻辑
package limiter

import (
	"strings"
	"thub/pkg/config"
	"thub/pkg/logger"
	"thub/pkg/redis"

	"github.com/gin-gonic/gin"
	limiterlib "github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// GetKeyIP 获取 Limitor 的 Key, IP
func GetKeyIP(c *gin.Context) string {
	return c.ClientIP()
}

// GetKeyRouteWithIP Limitor 的 key, 路由+IP, 针对单个路由做限流
func GetKeyRouteWithIP(c *gin.Context) string {
	return routeToKeyString(c.FullPath() + c.ClientIP())
}

// CheckRate 检测请求是否超额
func CheckRate(c *gin.Context, key, formatted string) (limiterlib.Context, error) {
	// 实例化依赖的 limiter 包的 limiter.Rate 对象
	var context limiterlib.Context
	rate, err := limiterlib.NewRateFromFormatted(formatted)
	if err != nil {
		logger.LogIf(err)
		return context, err
	}

	// 初始化存储 使用自定义的Redis
	store, err := sredis.NewStoreWithOptions(redis.Redis.Client, limiterlib.StoreOptions{
		Prefix: config.GetString("app.name") + ":limiter",
	})
	if err != nil {
		logger.LogIf(err)
		return context, err
	}

	// 组装实例
	limiterObj := limiterlib.New(store, rate)

	// 获取限流的结果
	if c.GetBool("limiter-once") {
		// Peek() 取结果, 不增加访问次数
		return limiterObj.Peek(c, key)
	} else {
		// 确保多个路由组里调用 LimitIP 进行限流时, 只增加一次访问次数
		c.Set("limiter-once", true)

		// Get() 取结果且增加访问次数
		return limiterObj.Get(c, key)
	}
}

// routeToKeyString 辅助方法, 将 URL 中的 / 格式改为 -
func routeToKeyString(routeName string) string {
	routeName = strings.ReplaceAll(routeName, "/", "-")
	routeName = strings.ReplaceAll(routeName, ":", "_")
	return routeName
}
