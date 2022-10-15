package requests

import (
	"thub/app/requests/validators"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// 手机重置密码参数验证器
type ResetByPhoneRequest struct {
	Phone      string `json:"phone,omitempty" valid:"phone"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
	Password   string `json:"password,omitempty" valid:"password"`
}

// 验证表单
func ResetByPhone(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"phone":       []string{"required", "digits:11"},
		"verify_code": []string{"required", "digits:6"},
		"password":    []string{"required", "min:6"},
	}

	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项",
			"digits:手机号长度必须为11位数字",
		},
		"verify_code": []string{
			"required:验证码为必填项",
			"digits:验证码长度为6",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度必须大于6",
		},
	}

	errs := validate(data, rules, messages)
	_data := data.(*ResetByPhoneRequest)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}

// 邮箱重置密码参数验证器
type ResetByEmailRequest struct {
	Email      string `json:"email,omitempty" valid:"email"`
	VerifyCode string `json:"verify_code,omitempty" valid:"verify_code"`
	Password   string `json:"password,omitempty" valid:"password"`
}

// 验证表单
func ResetByEmail(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"email":       []string{"required", "min:4", "max:30", "email"},
		"verify_code": []string{"required", "digits:6"},
		"password":    []string{"required", "min:6"},
	}

	messages := govalidator.MapData{
		"email": []string{
			"required:请填写邮箱",
			"min:邮箱长度最小为4",
			"max:邮箱长度最大为30",
			"email:请输入正确的邮箱格式",
		},
		"verify_code": []string{
			"required:请填写验证码",
			"digits:验证码长度为6",
		},
		"password": []string{
			"required:请填写密码",
			"min:密码长度应大于6",
		},
	}

	errs := validate(data, rules, messages)
	_data := data.(*ResetByEmailRequest)
	errs = validators.ValidateVerifyCode(_data.Email, _data.VerifyCode, errs)

	return errs
}
