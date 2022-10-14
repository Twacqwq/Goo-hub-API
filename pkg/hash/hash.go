// Package hash 哈希操作类
package hash

import (
	"thub/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// 使用 bcrypt 对密码加密
func BcryptHash(password string) string {
	// GenerateFromPassword 的第二个参数是 cost 值。建议大于 12，数值越大耗费时间越长
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger.LogIf(err)

	return string(bytes)
}

// 对比明文密码
func BcyptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 验证密码是否已Hash
func BcryptIsHashed(str string) bool {
	// Hashed len is 60
	return len(str) == 60
}
