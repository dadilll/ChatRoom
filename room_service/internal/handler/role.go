package handler

import (
	"net/http"
	"room_service/internal/models"
	"room_service/internal/service"
	"room_service/pkg/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RoleHandler struct {
	RoleService service.RoleService
	logger      logger.Logger
}

func NewRoleHandler(roleService service.RoleService, log logger.Logger) *RoleHandler {
	return &RoleHandler{
		RoleService: roleService,
		logger:      log,
	}
}

func (h *RoleHandler) CreateRole(c echo.Context) error {
	roomID := c.Param("room_id")

	var req models.Role
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	req.RoomID = roomID
	role, err := h.RoleService.CreateRole(req)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to create role", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create role"})
	}

	return c.JSON(http.StatusCreated, role)
}

func (h *RoleHandler) GetRole(c echo.Context) error {
	roleID := c.Param("role_id")

	role, err := h.RoleService.GetRole(roleID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to get role", zap.Error(err))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "role not found"})
	}

	return c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) GetRoomRoles(c echo.Context) error {
	roomID := c.Param("room_id")

	roles, err := h.RoleService.GetRolesByRoom(roomID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to get room roles", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get room roles"})
	}

	return c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) UpdateRole(c echo.Context) error {
	roleID := c.Param("role_id")

	var req models.UpdateRole
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	role, err := h.RoleService.UpdateRole(roleID, req)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to update role", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update role"})
	}

	return c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) DeleteRole(c echo.Context) error {
	roleID := c.Param("role_id")

	err := h.RoleService.DeleteRole(roleID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to delete role", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete role"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RoleHandler) AssignRole(c echo.Context) error {
	roomID := c.Param("room_id")
	userID := c.Param("user_id")

	var req struct {
		RoleID string `json:"role_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := h.RoleService.AssignRole(roomID, userID, req.RoleID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to assign role", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to assign role"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RoleHandler) RemoveRole(c echo.Context) error {
	roomID := c.Param("room_id")
	userID := c.Param("user_id")

	err := h.RoleService.RemoveRole(roomID, userID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to remove role", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to remove role"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RoleHandler) GetUserRole(c echo.Context) error {
	roomID := c.Param("room_id")
	userID := c.Param("user_id")

	role, err := h.RoleService.GetUserRole(roomID, userID)
	if err != nil {
		h.logger.Error(c.Request().Context(), "failed to get user role", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user role"})
	}

	return c.JSON(http.StatusOK, role)
}
