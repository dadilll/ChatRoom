package middleware

import (
	"crypto/rsa"
	"message_service/pkg/jwt"
	"net/http"
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
