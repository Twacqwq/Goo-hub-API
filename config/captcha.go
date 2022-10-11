package config

import "thub/pkg/config"

func init() {
	config.Add("captcha", func() map[string]interface{} {
		return map[string]interface{}{
			// 高度
			"height": 80,
			// 宽度
			"width": 240,
			// 长度
			"length": 6,
			// 最大倾斜角度
			"maxskew": 0.7,
			// 颗粒混淆数
			"dotcount": 80,
			// 过期时间(min)
			"expire_time": 15,
			//debug 模式下的exp
			"debug_expire_time": 10080,
			// 非 production 环境, 使用此 key 可跳过验证, 方便调试
			"testing_key": "captcha_skip_test",
		}
	})
}
