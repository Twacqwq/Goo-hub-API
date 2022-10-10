package auth

import (
	"fmt"
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

	// 解析 JSON 请求
	if err := c.ShouldBindJSON(&request); err != nil {
		// 解析失败 返回422状态码和错误信息
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err.Error())
		return
	}

	//表单验证
	errs := requests.ValidateSignupPhoneExist(&request, c)
	fmt.Printf("%v\n", errs)
	if len(errs) > 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"errors": errs,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}
