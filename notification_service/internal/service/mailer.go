package service

import (
	"fmt"
	"log"
	"net/smtp"
	"time"

	"notification_service/internal/config"
	"notification_service/internal/models"
)

type Mailer interface {
	Send(msg models.EmailMessage) error
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

func (m *SMTPMailer) Send(msg models.EmailMessage) error {
	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	address := fmt.Sprintf("%s:%d", m.host, m.port)

	// Заголовки
	date := time.Now().Format(time.RFC1123Z)
	messageID := fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), "id", m.host)

	// HTML тело письма
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<style>
				body { font-family: Arial, sans-serif; background: #f9f9f9; padding: 20px; }
				.container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.05); }
				h1 { color: #333; }
				p { color: #555; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>%s</h1>
				<p>%s</p>
			</div>
		</body>
		</html>
	`, msg.Subject, msg.Body)

	// Финальное письмо
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Date: %s\r\n"+
			"Message-ID: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s\r\n",
		m.from,
		msg.To,
		msg.Subject,
		date,
		messageID,
		htmlBody,
	))

	log.Printf("Отправка письма от %s к %s", m.from, msg.To)

	return smtp.SendMail(address, auth, m.from, []string{msg.To}, message)
}
