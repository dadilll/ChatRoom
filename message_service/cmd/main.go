package main

import (
	"context"
	"message_service/internal/config"
	kafka "message_service/internal/kafka/producer"
	"message_service/internal/server"
	"message_service/internal/service"
	"message_service/internal/storage"
	"message_service/pkg/db/postgres"
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

	db, err := postgres.New(cfg.ConfigPostgres)
	if err != nil {
		Logger.Error(ctx, "postgres error: "+err.Error())
		return
	}

	migrator := postgres.NewPostgresMigrator(db.Db)
	if err := migrator.Up(); err != nil {
		Logger.Error(ctx, "Error applying migrations: "+err.Error())
		return
	}

	kafkaWriter := kafka.NewWriterFromConfig(cfg.ConfigKafka)

	msgStorage := storage.NewPgMessageStorage(db.Db)
	msgService := service.NewMessageService(msgStorage)

	e := server.New(msgService, kafkaWriter)

	httpServer := server.Start(e, Logger, cfg.HTTPServerPort)
	server.WaitForShutdown(httpServer, Logger)
}
