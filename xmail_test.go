package xmail_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/coghost/xmail"
	"github.com/stretchr/testify/suite"
)

type XmailSuite struct {
	suite.Suite
}

// be remembered to set env variables to run tests
/**
export XMAIL_TO=xxx
export XMAIL_GF=xxx
export XMAIL_GP=xxx
export XMAIL_QF=xxx
export XMAIL_QP=xxx
**/

var (
	to = os.Getenv("XMAIL_TO")

	gCfg = &xmail.MailCfg{
		Port:     587,
		From:     os.Getenv("XMAIL_GF"),
		Password: os.Getenv("XMAIL_GP"),
		Host:     "smtp.gmail.com",
		Server:   xmail.GmailServer,
		To:       []string{to},
	}

	qCfg = &xmail.MailCfg{
		Port:     465,
		From:     os.Getenv("XMAIL_QF"),
		Password: os.Getenv("XMAIL_QP"),
		Host:     "smtp.exmail.qq.com",
		Server:   xmail.QQExmailServer,
		To:       []string{to},
	}
)

func init() {
	required := []string{
		"XMAIL_TO",
		"XMAIL_GF",
		"XMAIL_GP",
		"XMAIL_QF",
		"XMAIL_QP",
	}
	for _, v := range required {
		if _, b := os.LookupEnv(v); !b {
			panic(fmt.Sprintf("please set env %s first", v))
		}
	}
}

func TestXmail(t *testing.T) {
	suite.Run(t, new(XmailSuite))
}

func (s *XmailSuite) SetupSuite() {
}

func (s *XmailSuite) TearDownSuite() {
}

func (s *XmailSuite) Test_01() {
	qCfg.Alias = "Xmail"
	ms := xmail.GenMailService(qCfg)
	err := xmail.VerifyConfig(qCfg)
	s.Nil(err)

	subject, body := "Test with alias", "this is a test mail with alias Xmail"
	err = ms.Notify(subject, body)
	s.Nil(err)
}

type GmailSuite struct {
	suite.Suite
	cfg *xmail.MailCfg
}

func TestGmail(t *testing.T) {
	suite.Run(t, new(GmailSuite))
}

func (s *GmailSuite) SetupSuite() {
	s.cfg = gCfg
}

func (s *GmailSuite) TearDownSuite() {
}

func (s *GmailSuite) Test_00_no_alias() {
	ms := xmail.NewGmail(s.cfg)
	err := ms.Notify("Test no alias", "this is a test mail with no alias")
	s.Nil(err)
}

func (s *GmailSuite) Test_01_with_alias() {
	s.cfg.Alias = "Gmail"
	ms := xmail.NewGmail(s.cfg)
	err := ms.Notify("Test with alias", fmt.Sprintf("this is a test mail with alias %s", ms.Config.Alias))
	s.Nil(err)
}

type QQExmailSuite struct {
	suite.Suite
	cfg *xmail.MailCfg
}

func TestQQExmail(t *testing.T) {
	suite.Run(t, new(QQExmailSuite))
}

func (s *QQExmailSuite) SetupSuite() {
	qCfg.Alias = "QQExmail"
	s.cfg = qCfg
}

func (s *QQExmailSuite) TearDownSuite() {
}

func (s *QQExmailSuite) Test_00_no_alias() {
	ms := xmail.NewQQExmail(s.cfg)
	err := ms.Notify("Test no alias", "this is a test mail with no alias")
	s.Nil(err)
}

func (s *QQExmailSuite) Test_01_with_alias() {
	s.cfg.Alias = "QQExmail"
	ms := xmail.NewQQExmail(s.cfg)
	err := ms.Notify("Test with alias", fmt.Sprintf("this is a test mail with alias %s", ms.Config.Alias))
	s.Nil(err)
}
