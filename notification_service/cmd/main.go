package main

import (
	"log"
	"notification_service/internal/config"
	"notification_service/internal/consumer"
	"notification_service/internal/handler"
	mailer "notification_service/internal/service"
)

func main() {
	cfg, err := config.LoadMailerConfig()
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	mailer := mailer.NewSMTPMailer(*cfg)
	h := handler.NewNotificationHandler(mailer)

	go consumer.StartKafkaConsumer(cfg.KafkaBroker, cfg.KafkaTopic, h)

	select {} // блокируем main
}
