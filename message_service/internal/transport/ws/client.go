package ws

import (
	"context"
	"encoding/json"
	models "message_service/internal/DTO"
	producer "message_service/internal/kafka/producer"

	"github.com/gorilla/websocket"
	kafka "github.com/segmentio/kafka-go"
)

type Client struct {
	room        *Room
	conn        *websocket.Conn
	send        chan []byte
	KafkaWriter *producer.KafkaWriter
}

func NewClient(room *Room, conn *websocket.Conn, writer *producer.KafkaWriter) *Client {
	return &Client{
		room:        room,
		conn:        conn,
		send:        make(chan []byte, 256),
		KafkaWriter: writer,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg models.WSMessage[models.MessageRequest]
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		if msg.Event == models.EventMessageSend {
			resp, err := c.room.service.SendMessage(msg.Data)
			if err != nil {
				continue
			}

			kafkaMsg, _ := json.Marshal(models.WSMessage[models.MessageResponse]{
				Event: models.EventMessageReceive,
				Data:  resp,
			})

			c.KafkaWriter.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte(resp.RoomID),
				Value: kafkaMsg,
			})
		}
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
