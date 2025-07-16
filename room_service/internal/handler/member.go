package handler

import (
	"net/http"
	"room_service/internal/models"
	"room_service/internal/service"
	"room_service/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MemberHandler struct {
	service service.MemberService
	logger  logger.Logger
}

func NewMemberHandler(s service.MemberService, l logger.Logger) *MemberHandler {
	return &MemberHandler{
		service: s,
		logger:  l,
	}
}

func (h *MemberHandler) AddMember(c echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("room_id")

	userID, _ := c.Get("userID").(string)

	member := models.RoomMember{
		RoomID:   roomID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}

	err := h.service.AddMember(ctx, member)
	if err != nil {
		switch err.Error() {
		case "user already in room":
			return c.JSON(http.StatusConflict, map[string]string{"error": "already joined"})
		case "invite required":
			return c.JSON(http.StatusForbidden, map[string]string{"error": "invite required"})
		default:
			h.logger.Error(ctx, "failed to join room", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to join room"})
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (h *MemberHandler) ListMembers(c echo.Context) error {
	ctx := c.Request().Context()

	roomID := c.Param("room_id")
	members, err := h.service.ListMembers(ctx, roomID)
	if err != nil {
		h.logger.Error(ctx, "failed to list members", zap.String("room_id", roomID), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get members"})
	}

	return c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) RemoveMember(c echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("room_id")
	targetUserID := c.Param("user_id")

	requesterID, _ := c.Get("userID").(string)

	err := h.service.RemoveMember(ctx, roomID, targetUserID, requesterID)
	if err != nil {
		if err.Error() == "cannot remove other members" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "you can only leave the room yourself"})
		}

		h.logger.Error(ctx, "failed to leave room", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to leave room"})
	}

	return c.NoContent(http.StatusNoContent)
}
