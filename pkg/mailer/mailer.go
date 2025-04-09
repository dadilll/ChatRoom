package mailer

import (
	"fmt"
	"net/smtp"

	"service_auth/internal/config"
)

type Mailer interface {
	Send(to string, subject string, body string) error
}

type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPMailer(cfg config.MailerConfig) *SMTPMailer {
	return &SMTPMailer{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUser,
		password: cfg.SMTPPass,
		from:     cfg.SMTPFrom,
	}
}

func (m *SMTPMailer) Send(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s",
		m.from, to, subject, body))

	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	err := smtp.SendMail(addr, auth, m.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("ошибка при отправке письма: %w", err)
	}

	return nil
}
