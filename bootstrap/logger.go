package bootstrap

import (
	"thub/pkg/config"
	"thub/pkg/logger"
)

// 初始化 Logger
func SetupLogger() {
	logger.InitLogger(
		config.Get("log.filename"),
		config.GetInt("log.max_size"),
		config.GetInt("log.max_backup"),
		config.GetInt("log.max_age"),
		config.GetBool("log.compress"),
		config.Get("log.type"),
		config.Get("log.level"),
	)
}
