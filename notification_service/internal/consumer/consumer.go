package consumer

import (
	"context"
	"log"
	"notification_service/internal/handler"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func StartKafkaConsumer(brokerAddr, topic string, handler handler.Handler) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddr},
		Topic:       topic,
		GroupID:     "notification-group",
		StartOffset: kafka.LastOffset,
	})
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("Ошибка закрытия Kafka reader: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Получен сигнал завершения работы")
		cancel()
	}()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Завершение по сигналу.")
				return
			}
			log.Printf("Ошибка чтения из Kafka: %v", err)
			continue
		}

		log.Printf("Получено сообщение: %s", string(m.Value))

		if err := handler.Handle(m.Value); err != nil {
			log.Printf("Ошибка обработки уведомления: %v", err)
		}
	}
}
