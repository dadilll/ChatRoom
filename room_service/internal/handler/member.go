package handler

import (
	"net/http"
	"room_service/internal/models"
	"room_service/internal/service"
	"room_service/pkg/logger"
	"room_service/pkg/validator"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MemberHandler struct {
	service   service.MemberService
	logger    logger.Logger
	validator *validator.CustomValidator
}

func NewMemberHandler(s service.MemberService, l logger.Logger, v *validator.CustomValidator) *MemberHandler {
	return &MemberHandler{
		service:   s,
		logger:    l,
		validator: v,
	}
}

func (h *MemberHandler) AddMember(c echo.Context) error {
	req := models.AddMemberRequest{
		RoomID: c.Param("room_id"),
	}

	if errs := h.validator.Validate(&req); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	member := models.RoomMember{
		RoomID:   req.RoomID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}

	err := h.service.AddMember(c.Request().Context(), member)
	if err != nil {
		switch err.Error() {
		case "user already in room":
			return c.JSON(http.StatusConflict, map[string]string{"error": "already joined"})
		case "invite required":
			return c.JSON(http.StatusForbidden, map[string]string{"error": "invite required"})
		default:
			h.logger.Error(c.Request().Context(), "failed to join room", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to join room"})
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (h *MemberHandler) ListMembers(c echo.Context) error {
	req := models.ListMembersRequest{
		RoomID: c.Param("room_id"),
	}

	if errs := h.validator.Validate(&req); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	members, err := h.service.ListMembers(c.Request().Context(), req.RoomID, userID)
	if err != nil {
		if err.Error() == "forbidden: not a member of the room" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, members)
}

func (h *MemberHandler) RemoveMember(c echo.Context) error {
	req := models.RemoveMemberRequest{
		RoomID: c.Param("room_id"),
		UserID: c.Param("user_id"),
	}

	if errs := h.validator.Validate(&req); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	requesterID, ok := c.Get("userID").(string)
	if !ok || requesterID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	err := h.service.RemoveMember(c.Request().Context(), req.RoomID, req.UserID, requesterID)
	if err != nil {
		if err.Error() == "cannot remove other members" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "you can only leave the room yourself"})
		}

		h.logger.Error(c.Request().Context(), "failed to leave room", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to leave room"})
	}

	return c.NoContent(http.StatusNoContent)
}
