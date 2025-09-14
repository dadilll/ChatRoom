package router

import (
	"message_service/internal/handler"
	"message_service/internal/transport/ws"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	hub := ws.NewHub()

	e.GET("/ws/:roomID", handler.WebSocketHandler(hub))
}
