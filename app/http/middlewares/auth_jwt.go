// Package 认证中间件
package middlewares

import (
	"thub/app/models/user"
	"thub/pkg/jwt"
	"thub/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.NewJWT().ParseToken(c)
		if err != nil {
			// 没有权限
			response.Unauthorized(c, "没有权限访问")
			return
		}

		userModel := user.Get(claims.UserID)
		if userModel.ID == 0 {
			response.Unauthorized(c, "找不到该用户")
			return
		}

		// 将用户信息存入 gin.context 里  后续 auth 包从 context 拿数据
		c.Set("current_user_id", userModel.ID)
		c.Set("current_user_name", userModel.Name)
		c.Set("current_user", userModel)

		c.Next()
	}
}
