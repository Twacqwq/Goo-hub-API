package requests

import (
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
