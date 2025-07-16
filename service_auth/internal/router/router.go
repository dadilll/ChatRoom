package router

import (
	"crypto/rsa"
	"service_auth/internal/handler"
	middleware "service_auth/internal/middelware"
	"service_auth/internal/service"
	"service_auth/internal/storage"
	"service_auth/pkg/kafka"
	"service_auth/pkg/logger"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *sqlx.DB, Logger logger.Logger, privateKey *rsa.PrivateKey, rdb *redis.Client, kafkaWriter *kafka.KafkaWriter, publicKey *rsa.PublicKey) {
	userStorage := storage.NewUserStorage(db.DB)
	userStorageRedis := storage.NewRedisStorage(rdb)

	userService := service.NewUserService(userStorage, userStorageRedis, Logger, privateKey, kafkaWriter)
	userHandler := handler.NewUserHandler(userService, Logger)

	e.POST("api/v1/auth/login", userHandler.Login)
	e.POST("api/v1/auth/register", userHandler.Register)

	authGroup := e.Group("/api/v1/auth", middleware.AuthMiddleware(publicKey))

	authGroup.PUT("/profile/update", userHandler.ProfileUpdate)
	authGroup.GET("/profile", userHandler.GetProfile)
	authGroup.POST("/password/change", userHandler.ChangePassword)
	authGroup.POST("/email/confirm/send", userHandler.SendCodeEmail)
	authGroup.POST("/email/confirm/code", userHandler.ConfirmEmail)
	authGroup.POST("/email/change", userHandler.ChangeEmail)
	//authGroup.GET("/profile/{ID}", userHandler.GetProfileId)
}
