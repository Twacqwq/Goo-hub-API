package routes

import (
	"thub/app/http/controllers/api/v1/auth"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)

			// 判断手机是否注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)

			// 判断邮箱是否注册
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
		}
	}
}
