package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type ConfigKafka struct {
	KafkaBrokers []string `env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	Topic        string   `env:"KAFKA_TOPIC" env-default:"messages"`
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
	return NewWriter(cfg.KafkaBrokers, cfg.Topic)
}

func (w *KafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return w.writer.WriteMessages(ctx, msgs...)
}

func (w *KafkaWriter) Close() error {
	return w.writer.Close()
}
