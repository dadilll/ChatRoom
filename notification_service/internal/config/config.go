package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type MailerConfig struct {
	SMTPHost    string `env:"SMTP_HOST"`
	SMTPPort    int    `env:"SMTP_PORT"`
	SMTPUser    string `env:"SMTP_USERNAME"`
	SMTPPass    string `env:"SMTP_PASSWORD"`
	SMTPFrom    string `env:"SMTP_FROM"`
	KafkaBroker string `env:"KAFKA_BROKERS"`
	KafkaTopic  string `env:"KAFKA_EMAIL_TOPIC"`
}

func LoadMailerConfig() (*MailerConfig, error) {
	var cfg MailerConfig
	err := cleanenv.ReadConfig("conf/config.env", &cfg)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения конфигурации почты: %v", err)
	}
	return &cfg, nil
}
