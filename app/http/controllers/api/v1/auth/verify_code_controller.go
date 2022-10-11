package auth

import (
	v1 "thub/app/http/controllers/api/v1"
	"thub/pkg/captcha"
	"thub/pkg/logger"
	"thub/pkg/response"

	"github.com/gin-gonic/gin"
)

// 用户控制器
type VerifyCodeController struct {
	v1.BaseAPIController
}

func (vc *VerifyCodeController) ShowCaptcha(c *gin.Context) {
	// 生成验证码
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	// 记录错误日志
	logger.LogIf(err)
	// 返回数据
	response.JSON(c, gin.H{
		"captcha_id":    id,
		"captcha_image": b64s,
	})
}
