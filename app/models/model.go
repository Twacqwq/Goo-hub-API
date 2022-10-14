// Package models 模型通用属性和方法
package models

import (
	"time"

	"github.com/spf13/cast"
)

// 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autiIncrement;" json:"id,omitempty"` // omitempty 当字段为空, 序列化省略
}

// 获取 ID 的字符串格式
func (a BaseModel) GetStringID() string {
	return cast.ToString(a.ID)
}

// 时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"column:created_at;index" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;index" json:"updated_at,omitempty"`
}
