package config

import (
	"fmt"
	kafka "message_service/internal/kafka/producer"
	"message_service/pkg/db/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.ConfigPostgres
	kafka.ConfigKafka

	HTTPServerPort int    `env:"HTTP_SERVER_PORT" env-default:"8080"`
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
