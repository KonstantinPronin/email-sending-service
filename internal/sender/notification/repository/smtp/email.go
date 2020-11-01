package smtp

import (
	"crypto/tls"
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

	if err := e.sendMail(addr, e.conf.Host, e.conf.Login, notif.To, auth, msg); err != nil {
		e.logger.Error(fmt.Sprintf("Sending to smtp error: %s", err.Error()))
		return err
	}

	notif.SentStatus = true
	e.logger.Info(fmt.Sprintf("Message %s was sent to smtp", notif.ID))

	return nil
}

func (e *Email) sendMail(addr, host, from string, to []string, auth smtp.Auth, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err = conn.Close(); err != nil {
			e.logger.Error(fmt.Sprintf("closing resources error: %s", err.Error()))
		}
	}()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	for _, val := range to {
		if err = client.Rcpt(val); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	defer func() {
		if err = w.Close(); err != nil {
			e.logger.Error(fmt.Sprintf("closing resources error: %s", err.Error()))
		}
	}()

	if _, err = w.Write(msg); err != nil {
		return err
	}

	return client.Quit()
}

func NewEmail(
	conf *infrastructure.SmtpConf,
	logger *zap.Logger) internal.Transmitter {
	return &Email{conf: conf, logger: logger}
}
