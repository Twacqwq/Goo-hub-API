package requests

import (
	"thub/app/requests/validators"

	"github.com/gin-gonic/gin"

	"github.com/thedevsaddam/govalidator"
)

// 验证实体
type VerifyCodePhoneRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Phone         string `json:"phone,omitempty" valid:"phone"`
}

type VerifyCodeEmailRequest struct {
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
	Email         string `json:"email,omitempty" valid:"email"`
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
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}

// 验证表单
func VerifyCodeEmail(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
		"email":          []string{"required", "min:4", "max:30", "email"},
	}

	messages := govalidator.MapData{
		"captcha_id": []string{
			"required:图片的验证码为必填项",
		},
		"captcha_answer": []string{
			"required:图片验证码答案为必填项",
			"digits:图片验证码必须为6位数字",
		},
		"email": []string{
			"required:邮箱为必填项",
			"min:邮箱最小长度为4",
			"max:邮箱最大长度为30",
		},
	}

	errs := validate(data, rules, messages)

	// 图片验证码
	_data := data.(*VerifyCodeEmailRequest)
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}
