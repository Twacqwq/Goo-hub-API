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
			// 手机注册
			authGroup.POST("/signup/using-phone", suc.SignupUsingPhone)
			// 邮箱注册
			authGroup.POST("/signup/using-email", suc.SignupUsingEmail)

			vcc := new(auth.VerifyCodeController)

			// 显示验证码
			authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)
			// 发送手机验证码
			authGroup.POST("/verify-codes/phone", vcc.SendUsingPhone)
			// 发送邮箱验证码
			authGroup.POST("/verify-codes/email", vcc.SendUsingEmail)

			lgc := new(auth.LoginController)

			// 手机号+验证码登录
			authGroup.POST("/login/using-phone", lgc.LoginByPhone)
			// [手机号|用户名|邮箱] 登录
			authGroup.POST("/login/using-password", lgc.LoginByPassword)
			// 刷新令牌
			authGroup.POST("/login/refresh-token", lgc.RefreshToken)
		}
	}
}
