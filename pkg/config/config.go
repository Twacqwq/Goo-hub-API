package config

import (
	"os"
	"thub/pkg/helpers"

	"github.com/spf13/cast"
	viperlib "github.com/spf13/viper"
)

var viper *viperlib.Viper

// 动态加载配置信息
type ConfigFunc func() map[string]interface{}

// 先加载此map
var ConfigFuncs map[string]ConfigFunc

func init() {
	// 初始化
	viper = viperlib.New()

	// 配置类型
	viper.SetConfigType("env")

	// 环境变量配置文件查找的路径
	viper.AddConfigPath(".")

	// 设置前缀 区分系统变量
	viper.SetEnvPrefix("appenv")

	// 读取环境变量
	viper.AutomaticEnv()

	ConfigFuncs = make(map[string]ConfigFunc)
}

func InitConfig(env string) {
	// 加载环境变量
	loadEnv(env)
	// 注册配置信息
	loadConfig()
}

func loadConfig() {
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

func loadEnv(envSuffix string) {

	// 默认加载.env文件, 如果有传参 --env=name 加载.env.name文件
	envPath := ".env"
	if len(envSuffix) > 0 {
		filepath := ".env." + envSuffix
		if _, err := os.Stat(filepath); err == nil {
			envPath = filepath
		}
	}

	// 加载 .env
	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 监控 .env 文件 动态更新
	viper.WatchConfig()
}

// Env 读取环境变量, 支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return internalGet(envName, defaultValue...)
	}
	return internalGet(envName)
}

// 新增配置项
func Add(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

// Get 获取配置项
// 第一个参数 path 允许使用点式获取，如：app.name
// 第二个参数允许传参默认值
func Get(path string, defaultValue ...interface{}) string {
	return GetString(path, defaultValue...)
}

func internalGet(path string, defaultValue ...interface{}) interface{} {
	// config 或者环境变量不存在的情况
	if !viper.IsSet(path) || helpers.Empty(viper.Get(path)) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(path)
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString 获取结构数据
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}
