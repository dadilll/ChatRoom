package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ConfigPostgres struct {
	UserName string `env:"POSTGRES_USER" env-default:"postgres_user"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"P0stgr3sS3cur3"`
	Host     string `env:"POSTGRES_HOST" env-default:"db"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	DbName   string `env:"POSTGRES_DB" env-default:"authservice"`
}

type DB struct {
	Db *sqlx.DB
}

func New(config ConfigPostgres) (*DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", config.UserName, config.Password, config.DbName, config.Host, config.Port)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Db: db}, nil
}
