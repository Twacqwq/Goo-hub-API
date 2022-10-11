package auth

import (
	"net/http"
	v1 "thub/app/http/controllers/api/v1"
	"thub/app/models/user"
	"thub/app/requests"

	"github.com/gin-gonic/gin"
)

// 注册控制器
type SignupController struct {
	v1.BaseAPIController
}

// 检测手机号是否被注册
func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	//初始化请求对象
	request := requests.SignupPhoneExistRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateSignupPhoneExist); !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}

func (sc *SignupController) IsEmailExist(c *gin.Context) {
	request := requests.SignupEmailExistRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateSignupEmailExist); !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsEmailExist(request.Email),
	})
}
