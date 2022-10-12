package xmail

import (
	"errors"
)

var (
	ErrorNoServer    = errors.New("server is missing")
	ErrorNoRecipient = errors.New("recipients is missing")
)

const (
	GmailServer    = "gmail"
	QQExmailServer = "exmail"
)

type IMail interface {
	SendMail(subject, body string) error
}

type MailService struct {
	Config *MailCfg
}

type MailCfg struct {
	Port     int
	From     string
	Alias    string
	Password string
	Host     string
	// recipients
	To []string

	// for now only support exmail/gmail
	Server string
}

func NewMailService(cfg *MailCfg) *MailService {
	return &MailService{
		Config: cfg,
	}
}

func GenMailService(cfg *MailCfg) IMail {
	var ms IMail
	if cfg.Server == GmailServer {
		ms = NewGmail(cfg)
	} else {
		ms = NewQQExmail(cfg)
	}
	return ms
}

func SendMail(ms IMail, subject, body string) (err error) {
	return ms.SendMail(subject, body)
}

func VerifyConfig(cfg *MailCfg) error {
	if cfg.Host == "" {
		return ErrorNoServer
	}
	if len(cfg.To) == 0 {
		return ErrorNoRecipient
	}
	return nil
}
