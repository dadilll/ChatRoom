package server

import (
	"context"
	"fmt"
	"message_service/internal/router"
	"message_service/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()
	router.SetupRoutes(e)
	return e
}

func Start(server *echo.Echo, logger logger.Logger, port int) *http.Server {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: server,
	}

	go func() {
		logger.Info(context.Background(), fmt.Sprintf("Starting server on port :%d", port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "Failed to start server: "+err.Error())
		}
	}()

	return httpServer
}

func WaitForShutdown(httpServer *http.Server, logger logger.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info(context.Background(), "Received shutdown signal, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "Server Shutdown Failed: "+err.Error())
		return
	}

	logger.Info(context.Background(), "Server stopped gracefully")
}
