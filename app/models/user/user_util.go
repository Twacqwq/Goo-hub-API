package user

import "thub/pkg/database"

// 判断 Email 是否被注册
func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// 判断 Phone 是否被注册
func IsPhoneExist(phone string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

// GetByMulti 通过 手机号/Email/用户名 来获取用户
func GetByMulti(loginID string) (userModel User) {
	database.DB.
		Where("phone = ?", loginID).
		Or("email = ?", loginID).
		Or("name = ?", loginID).
		First(&userModel)
	return
}

// GetByPhone 通过手机号来获取用户
func GetByPhone(phone string) (userModel User) {
	database.DB.Where("phone = ?", phone).First(&userModel)
	return
}

// Get 通过ID获取用户
func Get(idstr string) (userModel User) {
	database.DB.Where("id", idstr).First(&userModel)
	return
}
