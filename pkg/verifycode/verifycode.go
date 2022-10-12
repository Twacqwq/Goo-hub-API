package verifycode

import (
	"strings"
	"sync"
	"thub/pkg/app"
	"thub/pkg/config"
	"thub/pkg/helpers"
	"thub/pkg/logger"
	"thub/pkg/redis"
	"thub/pkg/sms"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once
var verifyCode *VerifyCode

func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		verifyCode = &VerifyCode{
			Store: &RedisStore{
				RedisClient: redis.Redis,
				KeyPrefix:   config.GetString("app.name") + ":verifycode:",
			},
		}
	})
	return verifyCode
}

func (vs *VerifyCode) SendSMS(phone string) bool {
	// 生成验证码
	code := vs.generateVerifyCode(phone)

	// 方便本地调试
	if !app.IsProduction() && strings.HasPrefix(phone, config.GetString("verifycode.debug_phone_prefix")) {
		return true
	}

	// 发送短信
	return sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.aliyun.template_code"),
		Data:     map[string]string{"code": code},
	})

}

// 检查用户提交的验证码是否正确
func (vc *VerifyCode) CheckAnswer(key, answer string) bool {
	logger.DebugJSON("验证码", "检查验证码", map[string]string{key: answer})

	// 方便开发，在非生产环境下，具备特殊前缀的手机号和 Email后缀，会直接验证成功
	if !app.IsProduction() &&
		(strings.HasSuffix(key, config.GetString("verifycode.debug_email_suffix")) ||
			strings.HasPrefix(key, config.GetString("verifycode.debug_phone_prefix"))) {
		return true
	}

	// 这个false在线上环境要改为true
	return vc.Store.Verify(key, answer, false)
}

// 生成验证码并保存在Redis中
func (vc *VerifyCode) generateVerifyCode(key string) string {
	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))

	//为方便开发，本地调试使用固定验证码
	if app.IsLocal() {
		code = config.GetString("verifycode.debug_code")
	}

	logger.DebugJSON("验证码", "生成验证码", map[string]string{key: code})

	vc.Store.Set(key, code)
	return code
}
