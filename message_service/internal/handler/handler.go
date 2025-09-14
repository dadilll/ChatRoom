package handler

import (
	"net/http"

	"message_service/internal/transport/ws"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(hub *ws.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("roomID")
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		room := hub.GetRoom(roomID)
		client := ws.NewClient(room, conn)
		room.Register <- client

		go client.ReadPump()
		go client.WritePump()
		return nil
	}
}
