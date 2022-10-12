package auth

import (
	v1 "thub/app/http/controllers/api/v1"
	"thub/app/requests"
	"thub/pkg/captcha"
	"thub/pkg/logger"
	"thub/pkg/response"
	"thub/pkg/verifycode"

	"github.com/gin-gonic/gin"
)

// 用户控制器
type VerifyCodeController struct {
	v1.BaseAPIController
}

// 生成验证码
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

// 发送手机验证码
func (vc *VerifyCodeController) SendUsingPhone(c *gin.Context) {
	// 验证表单
	request := requests.VerifyCodePhoneRequest{}
	if ok := requests.Validate(c, &request, requests.VerifyCodePhone); !ok {
		return
	}

	// 发送SMS
	if ok := verifycode.NewVerifyCode().SendSMS(request.Phone); !ok {
		response.Abort500(c, "发送短信失败~")
	}
	response.Success(c)
}

// 发送邮箱验证码
func (vc *VerifyCodeController) SendUsingEmail(c *gin.Context) {
	request := requests.VerifyCodeEmailRequest{}
	if ok := requests.Validate(c, &request, requests.VerifyCodeEmail); !ok {
		return
	}

	// 发送 Email
	if err := verifycode.NewVerifyCode().SendEmail(request.Email); err != nil {
		response.Abort500(c, "发送邮箱验证码失败")
	} else {
		response.Success(c)
	}
}
