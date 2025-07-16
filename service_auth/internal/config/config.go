package config

import (
	"fmt"
	postgres "service_auth/pkg/db/postgres"
	redis "service_auth/pkg/db/redis"
	"service_auth/pkg/kafka"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.ConfigPostgres
	redis.ConfigRedis
	kafka.ConfigKafka

	HTTPServerPort int    `env:"HTTP_SERVER_PORT" env-default:"8080"`
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH" env-default:"key/private.pem"`
	PublicKeyPath  string `env:"PUBLIC_KEY" env-default:"key/public.pem"`
}

func New() *Config {
	cfg := Config{}
	err := cleanenv.ReadConfig("conf/conf.env", &cfg)
	if err != nil {
		fmt.Printf("Ошибка чтения конфигурации: %v\n", err)
		return nil
	}

	fmt.Printf("Загружена конфигурация: %+v\n", cfg)
	return &cfg
}
