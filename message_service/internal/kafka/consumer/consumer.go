package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	models "message_service/internal/DTO"
	"message_service/internal/transport/ws"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader *kafka.Reader
	Hub    *ws.Hub
}

func NewConsumer(reader *kafka.Reader, hub *ws.Hub) *Consumer {
	return &Consumer{
		Reader: reader,
		Hub:    hub,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		m, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Kafka read error:", err)
			continue
		}

		var msg models.WSMessage[models.MessageResponse]
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			fmt.Println("Failed to unmarshal Kafka message:", err)
			continue
		}

		room := c.Hub.GetRoom(msg.Data.RoomID)
		room.BroadcastMessage(m.Value)
	}
}
