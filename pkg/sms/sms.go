// Package sms 发送短信
package sms

import (
	"sync"
	"thub/pkg/config"
)

// 短信的结构体
type Message struct {
	Template string
	Data     map[string]string
	Content  string
}

// 操作结构体
type SMS struct {
	Driver Driver
}

// 确保全局唯一
var once sync.Once

var internalSMS *SMS

func NewSMS() *SMS {
	once.Do(func() {
		internalSMS = &SMS{
			Driver: &Aliyun{},
		}
	})
	return internalSMS
}

func (sms *SMS) Send(phone string, message Message) bool {
	return sms.Driver.Send(phone, message, config.GetStringMapString("sms.aliyun"))
}
