package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	postgres "room_service/pkg/db/postgres"
)

type Config struct {
	postgres.ConfigPostgres

	HTTPServerPort int    `env:"HTTP_SERVER_PORT" env-default:"8080"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH" env-default:"key/public.pem"`
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
