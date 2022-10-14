package requests

import (
	"thub/app/requests/validators"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// 手机+验证码登录验证器
type LoginByPhoneRequest struct {
	Phone      string `json:"phone,omitempty" valid:"phone"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
}

// LoginByPhone 验证表单, 返回长度等于零即通过
func LoginByPhone(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"phone":       []string{"required", "digits:11"},
		"verify_code": []string{"required", "digits:6"},
	}
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项",
			"digits:手机号长度为11位",
		},
		"verify_code": []string{
			"required:验证码为必填项",
			"digits:验证码长度为6位",
		},
	}

	errs := validate(data, rules, messages)
	_data := data.(*LoginByPhoneRequest)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}

// 使用 [手机号|用户名|邮箱] + 密码登录验证器
type LoginByPasswordRequest struct {
	LoginID       string `json:"login_id,omitempty" valid:"login_id"`
	Password      string `json:"password,omitempty" valid:"password"`
	CaptchaID     string `json:"captcha_id,omitempty" valid:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer,omitempty" valid:"captcha_answer"`
}

// LoginByPassword 验证表单
func LoginByPassword(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"login_id":       []string{"required", "min:3"},
		"password":       []string{"required", "min:6"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{
		"login_id": []string{
			"required:账号为必填项",
			"min:账号长度最小为3",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度最小为6",
		},
		"captcha_id": []string{
			"required:验证码ID为必填项",
		},
		"captcha_answer": []string{
			"required:验证码错误",
			"digits:验证码长度为6",
		},
	}

	errs := validate(data, rules, messages)
	_data := data.(*LoginByPasswordRequest)
	errs = validators.ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}
