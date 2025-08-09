package main

import (
	"context"
	"message_service/internal/config"
	"message_service/pkg/logger"
)

const serviceName = "message_service"

func main() {
	ctx := context.Background()
	Logger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, Logger)

	cfg := config.New()
	if cfg == nil {
		Logger.Error(ctx, "ERROR: config is nil")
		return
	}
}
