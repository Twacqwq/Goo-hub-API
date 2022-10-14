// Package 游客访问
package middlewares

import (
	"thub/pkg/jwt"
	"thub/pkg/response"

	"github.com/gin-gonic/gin"
)

func GuestJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.GetHeader("Authorization")) > 0 {
			_, err := jwt.NewJWT().ParseToken(c)
			// 解析 token 成功, 说明已登录
			if err == nil {
				response.Unauthorized(c, "请使用游客身份访问")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
