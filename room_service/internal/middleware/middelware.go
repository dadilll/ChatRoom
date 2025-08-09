package middleware

import (
	"crypto/rsa"
	"net/http"
	"room_service/internal/service"
	"room_service/pkg/jwt"
	"room_service/pkg/logger"
	"strings"

	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(publicKey *rsa.PublicKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
			}

			authParts := strings.SplitN(authHeader, " ", 2)
			if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token format"})
			}

			tokenString := authParts[1]
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Empty token"})
			}

			userID, err := jwt.ParseJWT(tokenString, publicKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}

			if userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing user ID in token"})
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

			ctx := c.Request().Context()
			hasPermission, err := roleService.CheckPermission(ctx, roomID, userID, permission)
			if err != nil {
				log.Error(ctx, "failed to check permission", zap.Error(err))
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check permissions"})
			}

			if !hasPermission {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "insufficient permissions"})
			}

			return next(c)
		}
	}
}
