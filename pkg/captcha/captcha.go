package captcha

import (
	"sync"
	"thub/pkg/app"
	"thub/pkg/config"
	"thub/pkg/redis"

	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

// 确保全局唯一
var once sync.Once

var internalCaptcha *Captcha

// 实例化
func NewCaptcha() *Captcha {
	once.Do(func() {
		internalCaptcha = &Captcha{}
		// 配置存储器
		store := &RedisStore{
			RedisClient: redis.Redis,
			KeyPrefix:   config.GetString("app.name") + ":captcha:",
		}
		// 配置驱动
		driver := base64Captcha.NewDriverDigit(
			config.GetInt("captcha.height"),
			config.GetInt("captcha.width"),
			config.GetInt("captcha.length"),
			config.GetFloat64("captcha.maxskew"), // 数字的最大倾斜角度
			config.GetInt("captcha.dotcount"),    // 图片背景的混淆点数量
		)

		// 实例化 base64Captcha 并赋值给内部使用的 internalCaptcha 对象
		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, store)
	})
	return internalCaptcha
}

// 生成图片验证码
func (c *Captcha) GenerateCaptcha() (id, b64s string, err error) {
	return c.Base64Captcha.Generate()
}

// 验证验证码是否正确
func (c *Captcha) VarifyCaptcha(id, answer string) (match bool) {
	// 方便本地和 API 调试
	if !app.IsProduction() && id == config.GetString("captcha.testing_key") {
		return true
	}

	return c.Base64Captcha.Verify(id, answer, true)
}
