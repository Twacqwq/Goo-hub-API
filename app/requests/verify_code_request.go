package requests

import (
	"thub/pkg/captcha"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// 验证实体
type VerifyCodePhoneRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Phone         string `json:"phone,omitempty" valid:"phone"`
}

// 验证表单
func VerifyCodePhone(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
		"phone":          []string{"required", "digits:11"},
	}

	messages := govalidator.MapData{
		"captcha_id": []string{
			"required:图片的验证码为必填项",
		},
		"captcha_answer": []string{
			"required:图片验证码答案为必填项",
			"digits:图片验证码必须为6位数字",
		},
		"phone": []string{
			"required:手机号为必填项",
			"digits:手机号必须为11位数字",
		},
	}

	errs := validate(data, rules, messages)

	// 图片验证码
	_data := data.(*VerifyCodePhoneRequest)
	if ok := captcha.NewCaptcha().VarifyCaptcha(_data.CaptchaID, _data.CaptchaAnswer); !ok {
		errs["captcha_answer"] = append(errs["captcha_answer"], "图片验证码错误")
	}

	return errs
}
