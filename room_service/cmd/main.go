package main

import (
	"context"
	"room_service/internal/config"
	"room_service/internal/server"
	"room_service/pkg/db/postgres"
	"room_service/pkg/jwt"
	"room_service/pkg/logger"
)

const serviceName = "room_service"

func main() {
	ctx := context.Background()
	Logger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, Logger)

	cfg := config.New()
	if cfg == nil {
		Logger.Error(ctx, "ERROR: config is nil")
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

	publicKey, err := jwt.LoadPublicKey(cfg.PublicKeyPath)
	if err != nil {
		Logger.Error(ctx, "Error loaded private key: "+err.Error())
	}

	e := server.New(db.Db, Logger, publicKey)

	httpServer := server.Start(e, Logger, cfg.HTTPServerPort)

	server.WaitForShutdown(httpServer, Logger)
}
