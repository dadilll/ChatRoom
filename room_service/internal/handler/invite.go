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

type InviteHandler struct {
	inviteService service.InviteService
	logger        logger.Logger
	validator     *validator.CustomValidator
}

func NewInviteHandler(svc service.InviteService, log logger.Logger, val *validator.CustomValidator) *InviteHandler {
	return &InviteHandler{
		inviteService: svc,
		logger:        log,
		validator:     val,
	}
}

func (h *InviteHandler) CreateInvite(c echo.Context) error {
	var req models.CreateInviteRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if errs := h.validator.Validate(&req); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	sentByID, ok := c.Get("userID").(string)
	if !ok || sentByID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user context"})
	}

	invite, err := h.inviteService.NewInvite(c.Request().Context(), req.RoomID, req.InvitedID, sentByID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to create invite", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create invite"})
	}

	return c.JSON(http.StatusCreated, invite)
}

func (h *InviteHandler) GetUserInvites(c echo.Context) error {
	params := models.UserIDParam{
		UserID: c.Param("user_id"),
	}
	if errs := h.validator.Validate(&params); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	invites, err := h.inviteService.GetUserInvites(c.Request().Context(), params.UserID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to get user invites", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get invites"})
	}
	return c.JSON(http.StatusOK, invites)
}

func (h *InviteHandler) AcceptInvite(c echo.Context) error {
	params := models.InviteIDParam{
		InviteID: c.Param("invite_id"),
	}
	if errs := h.validator.Validate(&params); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	err := h.inviteService.AcceptInvite(c.Request().Context(), params.InviteID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to accept invite", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to accept invite"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *InviteHandler) DeclineInvite(c echo.Context) error {
	params := models.InviteIDParam{
		InviteID: c.Param("invite_id"),
	}
	if errs := h.validator.Validate(&params); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	err := h.inviteService.DeclineInvite(c.Request().Context(), params.InviteID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to decline invite", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to decline invite"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *InviteHandler) DeleteInvite(c echo.Context) error {
	params := models.InviteIDParam{
		InviteID: c.Param("invite_id"),
	}
	if errs := h.validator.Validate(&params); errs != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errs,
		})
	}

	err := h.inviteService.DeleteInvite(c.Request().Context(), params.InviteID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to delete invite", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete invite"})
	}
	return c.NoContent(http.StatusNoContent)
}
