package bootstrap

import (
	"fmt"
	"thub/pkg/config"
	"thub/pkg/redis"
)

func SetupRedis() {
	redis.ConnectRedis(
		fmt.Sprintf("%v:%v", config.Get("redis.host"), config.Get("redis.port")),
		config.GetString("redis.username"),
		config.GetString("redis.password"),
		config.GetInt("redis.database"),
	)
}
