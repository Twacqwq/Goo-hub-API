package requests

import (
	"thub/app/requests/validators"
	"thub/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// 验证函数类型
type ValidatorFunc func(interface{}, *gin.Context) map[string][]string

type SignupPhoneExistRequest struct {
	Phone string `json:"phone,omitempty" valid:"phone"`
}

type SignupEmailExistRequest struct {
	Email string `json:"email,omitempty" valid:"email"`
}

// [回调]
func ValidateSignupPhoneExist(data interface{}, c *gin.Context) map[string][]string {
	// 自定义验证规则
	rules := govalidator.MapData{
		"phone": []string{"required", "digits:11"},
	}

	// 自定义验证错误信息
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项, 参数名称 phone",
			"digits:手机号长度必须为 11 位数字",
		},
	}

	return validate(data, rules, messages)
}

// [回调]
func ValidateSignupEmailExist(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"email": []string{"required", "min:4", "max:30", "email"},
	}

	messages := govalidator.MapData{
		"email": []string{
			"required:Email为必填项",
			"min:Email长度大于4",
			"max:Email长度小于30",
			"email:请提供有效的Email地址",
		},
	}

	return validate(data, rules, messages)
}

// 注册接口验证器
type SignupUsingPhoneRequest struct {
	Name            string `json:"name,omitempty" valid:"name"`
	VerifyCode      string `json:"verify_code,omitempty" valid:"verify_code"`
	Phone           string `json:"phone,omitempty" valid:"phone"`
	Password        string `json:"password,omitempty" valid:"password"`
	PasswordConfirm string `json:"password_confirm,omitempty" valid:"password_confirm"`
}

// 验证手机注册参数
func SignupUsingPhone(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20", "not_exists:users,name"},
		"verify_code":      []string{"required", "digits:6"},
		"phone":            []string{"required", "digits:11", "not_exists:users,phone"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}

	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone",
			"digits:手机号长度必须为 11 位的数字",
		},
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
		"verify_code": []string{
			"required:验证码答案必填",
			"digits:验证码长度必须为 6 位的数字",
		},
	}

	errs := validate(data, rules, messages)

	_data := data.(*SignupUsingPhoneRequest)
	errs = validators.ValidatePasswordConfirm(_data.Password, _data.PasswordConfirm, errs)
	errs = validators.ValidateVerifyCode(_data.Phone, _data.VerifyCode, errs)

	return errs
}

// 邮箱注册验证器
type SignupUsingEmailRequest struct {
	Email           string `json:"email,omitempty" valid:"email"`
	VerifyCode      string `json:"verify_code,omitempty" valid:"verify_code"`
	Name            string `json:"name" valid:"name"`
	Password        string `json:"password,omitempty" valid:"password"`
	PasswordConfirm string `json:"password_confirm,omitempty" valid:"password_confirm"`
}

// 验证邮箱注册参数
func SignupUsingEmail(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"email":            []string{"required", "email", "min:4", "max:30", "not_exists:users,email"},
		"verify_code":      []string{"required", "digits:6"},
		"name":             []string{"required", "alpha_num", "between:3,20", "not_exists:users,name"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}

	messages := govalidator.MapData{
		"email": []string{
			"required:邮箱为必填项",
			"email:邮箱格式不正确",
			"min:长度需要大于4",
			"max:长度不可超过30",
			"not_exists:该Email已被占用",
		},
		"verify_code": []string{
			"required:验证码为必填项",
			"digits:验证码长度为6",
		},
		"name": []string{
			"required:名称为必填项",
			"alpha_num:用户名格式错误,只允许英文和数字",
			"between:用户名长度在3~20之间",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度过小",
		},
		"password_confirm": []string{
			"required:请再次输入密码",
		},
	}

	errs := validate(data, rules, messages)
	_data := data.(*SignupUsingEmailRequest)
	errs = validators.ValidatePasswordConfirm(_data.Password, _data.PasswordConfirm, errs)
	errs = validators.ValidateVerifyCode(_data.Email, _data.VerifyCode, errs)

	return errs
}

func Validate(c *gin.Context, obj interface{}, handler ValidatorFunc) bool {
	// 解析请求，支持JSON数据，表单，URL Query
	if err := c.ShouldBind(obj); err != nil {
		response.BadRequest(c, err, "请求解析错误, 请确认请求格式是否正确")
		return false
	}
	// 验证表单
	errs := handler(obj, c)
	if len(errs) > 0 {
		response.ValidationError(c, errs)
		return false
	}
	return true
}

// 验证封装
func validate(data interface{}, rules, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}
	return govalidator.New(opts).ValidateStruct()
}
