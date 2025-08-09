package handler

import (
	"net/http"
	"room_service/internal/models"
	"room_service/internal/service"
	"room_service/pkg/logger"
	"room_service/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RoomHandler struct {
	RoomService service.RoomService
	logger      logger.Logger
	validator   *validator.CustomValidator
}

func NewRoomHandler(roomService service.RoomService, log logger.Logger, val *validator.CustomValidator) *RoomHandler {
	return &RoomHandler{
		RoomService: roomService,
		logger:      log,
		validator:   val,
	}
}

func (h *RoomHandler) CreateRoom(c echo.Context) error {
	var req models.CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	if errs := h.validator.Validate(&req); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "fields": errs})
	}

	if req.Private == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "private field is required"})
	}

	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		h.logger.Error(c.Request().Context(), "userID not found in context during room creation")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	room := models.RoomFromCreateRequest(req, userID)

	resp, err := h.RoomService.NewRoom(c.Request().Context(), room)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to create room", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create room"})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *RoomHandler) GetRoom(c echo.Context) error {
	roomID := c.Param("room_id")

	if roomID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "room_id is required"})
	}

	room, err := h.RoomService.GetRoom(c.Request().Context(), roomID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to get room", zap.String("room_id", roomID), zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "room not found"})
	}

	return c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) UpdateRoom(c echo.Context) error {
	roomID := c.Param("room_id")
	if roomID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "room_id is required"})
	}

	var req models.UpdateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request format"})
	}

	if errs := h.validator.Validate(&req); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "fields": errs})
	}

	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	room, err := h.RoomService.GetRoom(c.Request().Context(), roomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "room not found"})
	}

	if room.OwnerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "you are not the owner of the room"})
	}

	if err := h.RoomService.UpdateRoom(c.Request().Context(), roomID, req); err != nil {
		h.logger.Error(c.Request().Context(), "failed to update room", zap.String("room_id", roomID), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update room"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RoomHandler) DeleteRoom(c echo.Context) error {
	roomID := c.Param("room_id")
	if roomID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "room_id is required"})
	}

	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	room, err := h.RoomService.GetRoom(c.Request().Context(), roomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "room not found"})
	}

	if room.OwnerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "you are not the owner of the room"})
	}

	if err := h.RoomService.DeleteRoom(c.Request().Context(), roomID); err != nil {
		h.logger.Error(c.Request().Context(), "failed to delete room", zap.String("room_id", roomID), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete room"})
	}

	return c.NoContent(http.StatusNoContent)
}
