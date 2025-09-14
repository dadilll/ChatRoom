package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type ConfigKafka struct {
	KafkaBrokers []string `env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	EmailTopic   string   `env:"KAFKA_EMAIL_TOPIC" env-default:"email-topic"`
}

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewWriter(brokers []string, topic string) *KafkaWriter {
	return &KafkaWriter{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func NewWriterFromConfig(cfg ConfigKafka) *KafkaWriter {
	return NewWriter(cfg.KafkaBrokers, cfg.EmailTopic)
}

func (w *KafkaWriter) WriteMessage(ctx context.Context, msg []byte) error {
	return w.writer.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
}

func (w *KafkaWriter) Close() error {
	return w.writer.Close()
}
