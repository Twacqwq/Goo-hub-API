package verifycode

import (
	"thub/pkg/app"
	"thub/pkg/config"
	"thub/pkg/redis"
	"time"
)

type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

func (s *RedisStore) Set(key, value string) bool {
	expireTime := time.Minute * time.Duration(config.GetInt64("verifycode.expire_time"))
	if app.IsLocal() {
		expireTime = time.Minute * time.Duration(config.GetFloat64("verifycode.debug_expire_time"))
	}
	return s.RedisClient.Set(s.KeyPrefix+key, value, expireTime)
}

func (s *RedisStore) Get(key string, clear bool) string {
	key = s.KeyPrefix + key
	val := s.RedisClient.Get(key)
	if clear {
		s.RedisClient.Del(key)
	}
	return val
}

func (s *RedisStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
