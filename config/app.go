// Package config 站点配置信息
package config

import "thub/pkg/config"

func init() {
	config.Add("app", func() map[string]interface{} {
		return map[string]interface{}{

			// 应用名称
			"name": config.Env("APP_NAME", "THub"),

			// 当前环境 local | stage | production | test
			"env": config.Env("APP_ENV", "production"),

			// 是否进入调试模式
			"debug": config.Env("APP_DEBUG", false),

			// 应用服务端口
			"port": config.Env("APP_PORT", "9999"),

			// 加密JWT
			"key": config.Env("APP_KEY", "33446a9dcf9ea060a0a6532b166da32f304af0de"),

			// 用以生成链接
			"url": config.Env("APP_URL", "http://localhost:9999"),

			// 设置时区
			"timezone": config.Env("TIMEZONE", "Asia/Shanghai"),
		}
	})
}
