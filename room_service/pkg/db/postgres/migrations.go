package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Migrator interface {
	Up() error
	Down() error
}

type PostgresMigrator struct {
	db *sqlx.DB
}

func NewPostgresMigrator(db *sqlx.DB) *PostgresMigrator {
	return &PostgresMigrator{db: db}
}

func (m *PostgresMigrator) Up() error {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Применение миграций
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (m *PostgresMigrator) Down() error {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Откат миграций
	err = migrator.Down()
	if err != nil {
		return fmt.Errorf("failed to revert migrations: %w", err)
	}

	return nil
}
