package auth

import (
	v1 "thub/app/http/controllers/api/v1"
	"thub/app/requests"
	"thub/pkg/auth"
	"thub/pkg/jwt"
	"thub/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	v1.BaseAPIController
}

// 手机+验证码登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
		return
	}

	// 尝试登录
	user, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		response.Error(c, err, "账号不存在")
	} else {
		// 登录成功
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
		response.JSON(c, gin.H{
			"message": "登录成功",
			"token":   token,
		})
	}
}

func (lc *LoginController) LoginByPassword(c *gin.Context) {
	request := requests.LoginByPasswordRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPassword); !ok {
		return
	}

	// 尝试登录
	user, err := auth.Attempt(request.LoginID, request.Password)
	if err != nil {
		response.Unauthorized(c, "账号不存在或密码错误")

	} else {
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
		response.JSON(c, gin.H{
			"message": "登录成功",
			"token":   token,
		})
	}
}
