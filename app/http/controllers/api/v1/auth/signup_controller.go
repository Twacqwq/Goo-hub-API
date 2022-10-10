package auth

import (
	"fmt"
	"net/http"
	v1 "thub/app/http/controllers/api/v1"
	"thub/app/models/user"

	"github.com/gin-gonic/gin"
)

// 注册控制器
type SignupController struct {
	v1.BaseAPIController
}

// 检测手机号是否被注册
func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	// 请求对象
	type PhoneExistRequest struct {
		Phone string `json:"phone"`
	}
	request := PhoneExistRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		// 解析失败 返回422状态码和错误信息
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		fmt.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}
