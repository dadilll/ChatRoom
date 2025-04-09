package middleware

import (
	"crypto/rsa"
	"net/http"
	"service_auth/pkg/jwt"
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
