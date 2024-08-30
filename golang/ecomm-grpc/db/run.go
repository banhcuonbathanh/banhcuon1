package db

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return pool, nil
}

func RunMigrations(databaseURL string) error {
    migrationsPath := "file://ecomm-grpc/db/migrations" // Relative path from the project root

	// migrationsPath1 := fmt.Sprintf("file://%s/db/migrations", os.Getwd())
    m, err := migrate.New(migrationsPath, databaseURL)
    if err != nil {
        return fmt.Errorf("failed to create migrate instance: %w", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    return nil
}

// func RunMigrations(databaseURL string) error {
//     migrationsPath := "file://ecomm-grpc/db/migrations"

//     m, err := migrate.New(migrationsPath, databaseURL)
//     if err != nil {
//         return fmt.Errorf("failed to create migrate instance: %w", err)
//     }
//     defer m.Close()

//     // Force the version to 7 (or the version you want to force)
//     if err := m.Force(7); err != nil {
//         return fmt.Errorf("failed to force migration version: %w", err)
//     }

//     // Now run the migrations
//     if err := m.Up(); err != nil && err != migrate.ErrNoChange {
//         return fmt.Errorf("failed to run migrations: %w", err)
//     }

//     return nil
// }