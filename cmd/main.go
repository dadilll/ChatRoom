package main

import (
	"context"
	"service_auth/internal/config"
	"service_auth/internal/server"
	"service_auth/pkg/db/postgres"
	"service_auth/pkg/db/redis"
	"service_auth/pkg/jwt"
	"service_auth/pkg/logger"
	"service_auth/pkg/mailer"
)

const serviceName = "Auth_service"

func main() {
	ctx := context.Background()
	Logger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, Logger)

	mailerCfg, err := config.LoadMailerConfig()
	if err != nil {
		Logger.Error(ctx, "Error loading mailer config: "+err.Error())
		return
	}

	cfg := config.New()
	if cfg == nil {
		Logger.Error(ctx, "ERROR: config is nil")
		return
	}

	rdb, err := redis.New(cfg.ConfigRedis)
	if err != nil {
		Logger.Error(ctx, "redis connection error: "+err.Error())
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

	privateKey, err := jwt.LoadPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		Logger.Error(ctx, "Error loaded private key: "+err.Error())
	}

	mailer := mailer.NewSMTPMailer(*mailerCfg)

	// инициализируем сервер
	e := server.New(db.Db, Logger, privateKey, rdb, mailer)

	httpServer := server.Start(e, Logger, cfg.HTTPServerPort)

	server.WaitForShutdown(httpServer, Logger)
}
