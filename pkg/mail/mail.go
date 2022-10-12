package mail

import (
	"sync"
	"thub/pkg/config"
)

type From struct {
	Address string
	Name    string
}

type Email struct {
	From    From
	To      []string
	Bcc     []string
	Cc      []string
	Subject string
	Text    []byte
	HTML    []byte
}

type Mailer struct {
	Driver Driver
}

var once sync.Once
var internalMailer *Mailer

func NewMailer() *Mailer {
	once.Do(func() {
		internalMailer = &Mailer{
			Driver: &SMTP{},
		}
	})
	return internalMailer
}

// 发送邮件验证码
func (im *Mailer) Send(email Email) bool {
	return im.Driver.Send(email, config.GetStringMapString("mail.smtp"))
}
