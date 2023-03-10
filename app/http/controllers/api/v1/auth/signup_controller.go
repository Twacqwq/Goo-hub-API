package auth

import (
	v1 "thub/app/http/controllers/api/v1"
	"thub/app/models/user"
	"thub/app/requests"
	"thub/pkg/jwt"
	"thub/pkg/response"

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
	response.JSON(c, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}

// 检测邮箱是否被注册
func (sc *SignupController) IsEmailExist(c *gin.Context) {
	request := requests.SignupEmailExistRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateSignupEmailExist); !ok {
		return
	}
	response.JSON(c, gin.H{
		"exist": user.IsEmailExist(request.Email),
	})
}

// 手机注册
func (sc *SignupController) SignupUsingPhone(c *gin.Context) {
	request := requests.SignupUsingPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingPhone); !ok {
		return
	}

	_user := user.User{
		Name:     request.Name,
		Phone:    request.Phone,
		Password: request.Password,
	}

	// 创建用户
	_user.Create()

	if _user.ID > 0 {
		token := jwt.NewJWT().IssueToken(_user.GetStringID(), _user.Name)
		response.CreatedJSON(c, gin.H{
			"token": token,
			"data":  _user,
		})
	} else {
		response.Abort500(c, "创建用户失败, 请稍后再试")
	}
}

// 邮箱注册
func (sc *SignupController) SignupUsingEmail(c *gin.Context) {
	request := requests.SignupUsingEmailRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingEmail); !ok {
		return
	}

	userModel := user.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	userModel.Create()

	if userModel.ID > 0 {
		token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Name)
		response.CreatedJSON(c, gin.H{
			"token": token,
			"data":  userModel,
		})
	} else {
		response.Abort500(c, "创建用户失败, 请稍后再试")
	}
}
