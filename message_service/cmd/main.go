package main

import (
	"context"
	"message_service/internal/config"
	"message_service/internal/server"
	"message_service/pkg/logger"
)

const serviceName = "message_service"

func main() {
	ctx := context.Background()
	Logger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, Logger)

	cfg := config.New()
	if cfg == nil {
		Logger.Error(ctx, "config is nil")
		return
	}

	e := server.New()

	httpServer := server.Start(e, Logger, cfg.HTTPServerPort)

	server.WaitForShutdown(httpServer, Logger)
}
