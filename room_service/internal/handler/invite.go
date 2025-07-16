package handler

import (
	"net/http"
	"room_service/internal/models"
	"room_service/internal/service"

	"github.com/labstack/echo/v4"
)

type InviteHandler struct {
	inviteService service.InviteService
}

func NewInviteHandler(svc service.InviteService) *InviteHandler {
	return &InviteHandler{
		inviteService: svc,
	}
}

func (h *InviteHandler) CreateInvite(c echo.Context) error {
	var req models.CreateInviteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	sentByID, ok := c.Get("userID").(string)
	if !ok || sentByID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user context")
	}

	invite, err := h.inviteService.NewInvite(c.Request().Context(), req.RoomID, req.InvitedID, sentByID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, invite)
}

func (h *InviteHandler) GetUserInvites(c echo.Context) error {
	userID := c.Param("user_id")
	invites, err := h.inviteService.GetUserInvites(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, invites)
}

func (h *InviteHandler) AcceptInvite(c echo.Context) error {
	id := c.Param("invite_id")
	err := h.inviteService.AcceptInvite(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *InviteHandler) DeclineInvite(c echo.Context) error {
	id := c.Param("invite_id")
	err := h.inviteService.DeclineInvite(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *InviteHandler) DeleteInvite(c echo.Context) error {
	id := c.Param("invite_id")
	err := h.inviteService.DeleteInvite(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
