package xmail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

type QQExmailService struct {
	MailService
}

func NewQQExmail(cfg *MailCfg) *QQExmailService {
	return &QQExmailService{
		*NewMailService(cfg),
	}
}

func (s *QQExmailService) Notify(subject, body string) (err error) {
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
	auth := smtp.PlainAuth("", s.Config.From, s.Config.Password, s.Config.Host)

	msg := s.createMessage(subject, body)
	return sendMailWithTLS(addr, auth, s.Config.From, s.Config.To, *msg)
}

func (s *QQExmailService) createMessage(subject, body string) *string {
	header := make(map[string]string)
	from := s.Config.From
	if s.Config.Alias != "" {
		from = fmt.Sprintf("%s<%s>", s.Config.Alias, s.Config.From)
	}
	header["From"] = from
	header["Subject"] = subject
	// format with html
	header["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	return &message
}

// sendMailWithTLS send email with tls
// refer: https://cloud.tencent.com/document/product/1288/65752
func sendMailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg string) (err error) {
	//create smtp client
	c, err := dial(addr)
	if err != nil {
		log.Println("Create smtp client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// dial return a smtp client
func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("tls.Dial Error:", err)
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
