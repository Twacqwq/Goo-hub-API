package validators

import (
	"errors"
	"fmt"
	"strings"
	"thub/pkg/database"

	"github.com/thedevsaddam/govalidator"
)

func init() {
	// govalidator 自定义验证器
	govalidator.AddCustomRule("not_exists", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		// 第一个参数为表名称 e.g users
		tableName := rng[0]
		// 第二个参数为字段名称 e.g email or phone
		dbField := rng[1]

		// 第三个参数 排除ID
		var exceptID string
		if len(rng) > 2 {
			exceptID = rng[2]
		}

		// 请求过来的数据
		requestValue := value.(string)

		// 拼接query
		query := database.DB.Table(tableName).Where(dbField+" = ?", requestValue)

		// 如果有第三个参数
		if len(exceptID) > 0 {
			query.Where("id != ?", exceptID)
		}

		var count int64
		query.Count(&count)

		if count != 0 {
			// 如果有自定义消息
			if message != "" {
				return errors.New(message)
			}
			// 默认的消息
			return fmt.Errorf("%v 也被占用", requestValue)
		}

		return nil
	})
}
