package smtp

import (
	"fmt"
	"github.com/KonstantinPronin/email-sending-service/internal"
	"github.com/KonstantinPronin/email-sending-service/pkg/infrastructure"
	"github.com/KonstantinPronin/email-sending-service/pkg/model"
	"go.uber.org/zap"
	"net/smtp"
	"strings"
)

type Email struct {
	conf   *infrastructure.SmtpConf
	logger *zap.Logger
}

func (e *Email) Transfer(notif *model.Notification) error {
	addr := strings.Join([]string{e.conf.Host, e.conf.Port}, ":")
	msg := []byte(fmt.Sprintf("To: %s \r\nSubject: %s\r\n\r\n%s\r\n",
		strings.Join(notif.To, ", "), notif.Subject, notif.Message))
	auth := smtp.PlainAuth("", e.conf.Login, e.conf.Password, e.conf.Host)

	if err := smtp.SendMail(addr, auth, e.conf.Login, notif.To, msg); err != nil {
		return err
	}

	notif.SentStatus = true
	e.logger.Info(fmt.Sprintf("Message %s was sent to smtp", notif.ID))

	return nil
}

func NewEmail(
	conf *infrastructure.SmtpConf,
	logger *zap.Logger) internal.Transmitter {
	return &Email{conf: conf, logger: logger}
}
