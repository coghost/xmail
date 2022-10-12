package xmail

import (
	"fmt"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
)

type GmailService struct {
	MailService
}

func NewGmail(cfg *MailCfg) *GmailService {
	return &GmailService{
		*NewMailService(cfg),
	}
}

func (s *GmailService) Notify(subject, body string) (err error) {
	client := &email.Email{
		From:    s.Config.From,
		To:      s.Config.To,
		Subject: subject,
		HTML:    []byte(body),
		Headers: textproto.MIMEHeader{},
	}
	if s.Config.Alias != "" {
		client.From = fmt.Sprintf("%s<%s>", s.Config.Alias, s.Config.From)
	}

	return client.Send(
		fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port),
		smtp.PlainAuth("", s.Config.From, s.Config.Password, s.Config.Host),
	)
}
