// Package models 模型通用属性和方法
package models

import "time"

// 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autiIncrement;" json:"id,omitempty"` // omitempty 当字段为空, 序列化省略
}

// 时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"column:created_at;index" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;index" json:"updated_at,omitempty"`
}
