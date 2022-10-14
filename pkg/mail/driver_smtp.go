package mail

import (
	"thub/pkg/logger"

	"github.com/spf13/cast"
	"gopkg.in/gomail.v2"
)

// SMTP 实现 mail.Driver interface
type SMTP struct{}

func (s *SMTP) Send(email Email, config map[string]string) bool {
	m := gomail.NewMessage()

	// 发送者
	m.SetHeader("From", m.FormatAddress(email.From.Address, email.From.Name))
	// 接收者
	m.SetHeader("To", email.To...)
	// m.SetAddressHeader("To", "example.com", "user")
	// 抄送者
	m.SetHeader("Cc", email.Cc...)
	// 邮件标题
	m.SetHeader("Subject", email.Subject)
	// 邮件正文
	m.SetBody("text/html", email.HTMLBody)
	// 邮件附件
	if len(email.File) > 0 {
		for _, files := range email.File {
			m.Attach(files.FilePath, gomail.Rename(files.FileName))
		}
	}

	logger.DebugJSON("发送邮件", "发送详情", email)

	d := gomail.NewDialer(config["host"], cast.ToInt(config["port"]), config["username"], config["password"])
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // 关闭TLS认证

	// Send
	if err := d.DialAndSend(m); err != nil {
		logger.ErrorString("发送邮件", "发送出错", err.Error())
		return false
	}

	logger.DebugString("发送邮件", "发送成功", "")
	return true
}
