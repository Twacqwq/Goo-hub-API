package user

import (
	"thub/app/models"
	"thub/pkg/database"
	"thub/pkg/hash"
)

type User struct {
	models.BaseModel
	Name     string `json:"name,omitempty"`
	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	models.CommonTimestampsField
}

// 插入一条数据
func (user *User) Create() {
	database.DB.Create(&user)
}

// 匹对密码是否正确
func (user *User) ComparePassword(_password string) bool {
	return hash.BcyptCheck(_password, user.Password)
}
