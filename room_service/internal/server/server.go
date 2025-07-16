package server

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"room_service/internal/router"
	"room_service/pkg/logger"
	"syscall"
	"time"
)

func New(db *sqlx.DB, Logger logger.Logger, publicKey *rsa.PublicKey) *echo.Echo {
	e := echo.New()
	router.SetupRoutes(e, Logger, db, publicKey)
	return e
}
func Start(server *echo.Echo, Logger logger.Logger, port int) *http.Server {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: server,
	}

	go func() {
		Logger.Info(context.Background(), fmt.Sprintf("Starting server on port :%d", port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Error(context.Background(), "Failed to start server: "+err.Error())
		}
	}()

	return httpServer
}
func WaitForShutdown(httpServer *http.Server, Logger logger.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	Logger.Info(context.Background(), "Received shutdown signal, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		Logger.Error(context.Background(), "Server Shutdown Failed: "+err.Error())
		return
	}

	Logger.Info(context.Background(), "Server stopped gracefully")
}
