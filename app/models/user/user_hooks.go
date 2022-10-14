package user

import (
	"thub/pkg/hash"

	"gorm.io/gorm"
)

// 密码加密钩子
func (userModel *User) BeforeSave(tx *gorm.DB) (err error) {
	if !hash.BcryptIsHashed(userModel.Password) {
		userModel.Password = hash.BcryptHash(userModel.Password)
	}

	return
}
