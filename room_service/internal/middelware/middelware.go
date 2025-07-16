package middleware

import (
	"context"
	"crypto/rsa"
	"go.uber.org/zap"
	"net/http"
	"room_service/internal/service"
	"room_service/pkg/jwt"
	"room_service/pkg/logger"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(publicKey *rsa.PublicKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token format"})
			}

			userID, err := jwt.ParseJWT(tokenString, publicKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}

			c.Set("userID", userID)
			return next(c)
		}
	}
}

func RequirePermission(roleService service.RoleService, log logger.Logger, permission int64) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roomID := c.Param("room_id")
			if roomID == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "room_id is required"})
			}

			userID, ok := c.Get("userID").(string)
			if !ok || userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			hasPermission, err := roleService.CheckPermission(roomID, userID, permission)
			if err != nil {
				log.Error(context.Background(), "failed to check permission", zap.Error(err))
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check permissions"})
			}

			if !hasPermission {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "insufficient permissions"})
			}

			return next(c)
		}
	}
}
