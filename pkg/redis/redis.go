// Package redis 工具包
package redis

import (
	"context"
	"sync"
	"thub/pkg/logger"
	"time"

	redis "github.com/go-redis/redis/v8"
)

// Redis 服务
type RedisClient struct {
	Client  *redis.Client
	Context context.Context
}

// 确保全局的 Redis 对象只实例一次
var once sync.Once

// Redis 全局 Redis
var Redis *RedisClient

// 连接 Redis
func ConnectRedis(addr, username, password string, db int) {
	once.Do(func() {
		Redis = NewClient(addr, username, password, db)
	})
}

// 实例化 Redis 连接
func NewClient(addr, username, password string, db int) *RedisClient {
	rds := &RedisClient{Context: context.Background()}
	rds.Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	// 测试连接
	err := rds.Ping()
	logger.LogIf(err)

	return rds
}

// 测试连接
func (rds RedisClient) Ping() error {
	_, err := rds.Client.Ping(rds.Context).Result()
	return err
}

// 存储key-value并设置exp
func (rds RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	if err := rds.Client.Set(rds.Context, key, value, expiration).Err(); err != nil {
		logger.ErrorString("Redis", "Set", err.Error())
		return false
	}
	return true
}

// 根据key获取value
func (rds RedisClient) Get(key string) string {
	result, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "Get", err.Error())
		}
		return ""
	}
	return result
}

// 判断key是否存在
func (rds RedisClient) Has(key string) bool {
	_, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "Has", err.Error())
		}
		return false
	}
	return true
}

// 删除keys
func (rds RedisClient) Del(keys ...string) bool {
	if err := rds.Client.Del(rds.Context, keys...).Err(); err != nil {
		logger.ErrorString("Redis", "Del", err.Error())
		return false
	}
	return true
}

// 清空数据库
func (rds RedisClient) FlushDB() bool {
	if err := rds.Client.FlushDB(rds.Context).Err(); err != nil {
		logger.ErrorString("Redis", "FlushDB", err.Error())
		return false
	}
	return true
}

// Increment 自增 当参数只有一个时,为key，其值增加1
// 两个参数时 第一个参数为key 第二个为增加的值 int64
func (rds RedisClient) Increment(params ...interface{}) bool {
	switch len(params) {
	case 1:
		key := params[0].(string)
		if err := rds.Client.Incr(rds.Context, key).Err(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false
		}
	case 2:
		key := params[0].(string)
		value := params[1].(int64)
		if err := rds.Client.IncrBy(rds.Context, key, value).Err(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false
		}
	default:
		logger.ErrorString("Redis", "Increment", "参数过多")
		return false
	}
	return true
}

// Decrement 当参数只有 1 个时，为 key，其值减去 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型
func (rds RedisClient) Decrement(params ...interface{}) bool {
	switch len(params) {
	case 1:
		key := params[0].(string)
		if err := rds.Client.Decr(rds.Context, key).Err(); err != nil {
			logger.ErrorString("Redis", "Decrement", err.Error())
			return false
		}
	case 2:
		key := params[0].(string)
		value := params[1].(int64)
		if err := rds.Client.DecrBy(rds.Context, key, value).Err(); err != nil {
			logger.ErrorString("redis", "Decrement", err.Error())
			return false
		}
	default:
		logger.ErrorString("Redis", "Decrement", "参数过多")
		return false
	}
	return true
}
