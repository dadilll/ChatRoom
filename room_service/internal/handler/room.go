package handler

import (
	"go.uber.org/zap"
	"net/http"
	"room_service/internal/models"
	"room_service/internal/service"
	"room_service/pkg/logger"

	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	RoomService service.RoomService
	logger      logger.Logger
}

func NewRoomHandler(roomService service.RoomService, log logger.Logger) *RoomHandler {
	return &RoomHandler{
		RoomService: roomService,
		logger:      log,
	}
}

func (h *RoomHandler) CreateRoom(c echo.Context) error {
	var req models.Room
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	userID, _ := c.Get("userID").(string)
	req.OwnerID = userID
	resp, err := h.RoomService.NewRoom(req)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to create room", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create room"})
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *RoomHandler) GetRoom(c echo.Context) error {
	id := c.Param("room_id")
	room, err := h.RoomService.GetRoom(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "room not found"})
	}
	return c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) UpdateRoom(c echo.Context) error {
	id := c.Param("room_id")
	var req models.UpdateRoom
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid data"})
	}
	err := h.RoomService.UpdateRoom(id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update room"})
	}
	return c.NoContent(http.StatusOK)
}

func (h *RoomHandler) DeleteRoom(c echo.Context) error {
	id := c.Param("room_id")
	err := h.RoomService.DeleteRoom(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete room"})
	}
	return c.NoContent(http.StatusOK)
}
