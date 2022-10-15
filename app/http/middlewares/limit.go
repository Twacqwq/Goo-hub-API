package middlewares

import (
	"net/http"
	"thub/pkg/app"
	"thub/pkg/limiter"
	"thub/pkg/logger"
	"thub/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// LimitIP 全局限流中间件，针对 IP 进行限流
// limit 为格式化字符串，如 "5-S" ，示例:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
func LimitIP(limit string) gin.HandlerFunc {
	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(c *gin.Context) {
		// 针对 IP 限流
		key := limiter.GetKeyIP(c)
		if ok := limitHandler(c, key, limit); !ok {
			return
		}

		c.Next()
	}
}

// LimitPerRoute 限流中间件, 用在单独的路由上
func LimitPerRoute(limit string) gin.HandlerFunc {
	if app.IsTesting() {
		limit = "1000000-H"
	}

	return func(c *gin.Context) {
		// 针对单个路由, 增加访问次数
		c.Set("limiter-once", false)
		// 针对IP + 路由进行限流
		key := limiter.GetKeyRouteWithIP(c)
		if ok := limitHandler(c, key, limit); !ok {
			return
		}

		c.Next()
	}
}

func limitHandler(c *gin.Context, key, limit string) bool {
	// 获取限额的情况
	rate, err := limiter.CheckRate(c, key, limit)
	if err != nil {
		logger.LogIf(err)
		response.Abort500(c)
		return false
	}

	// 设置Header信息
	c.Header("X-RateLimit-Limit", cast.ToString(rate.Limit))         // 最大访问次数
	c.Header("X-RateLimit-Remaining", cast.ToString(rate.Remaining)) // 剩余访问次数
	c.Header("X-RateLimit-Reset", cast.ToString(rate.Reset))         // 到达某个时间点, 访问次数重置为 X-RateLimit-Limit

	// 超额
	if rate.Reached {
		// 提示用户超额
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "请求接口太频繁",
		})
		return false
	}

	return true
}
