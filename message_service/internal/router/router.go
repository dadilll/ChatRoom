package router

import (
	"message_service/internal/handler"
	kafka "message_service/internal/kafka/producer"
	"message_service/internal/service"
	"message_service/internal/transport/ws"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, msgService service.MessageService, kafkaWriter *kafka.KafkaWriter) {
	hub := ws.NewHub(msgService)

	e.GET("/ws/:roomID", handler.WebSocketHandler(hub, kafkaWriter))
}
