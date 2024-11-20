package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

// func Connect(databaseURL string) (*pgxpool.Pool, error) {
// 	config, err := pgxpool.ParseConfig(databaseURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to parse database URL: %v", err)
// 	}

// 	pool, err := pgxpool.ConnectConfig(context.Background(), config)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to connect to database: %v", err)
// 	}

// 	return pool, nil
// }

// func RunMigrations(databaseURL string) error {
//     migrationsPath := "file://ecomm-grpc/db/migrations" // Relative path from the project root

// 	// migrationsPath1 := fmt.Sprintf("file://%s/db/migrations", os.Getwd())
//     m, err := migrate.New(migrationsPath, databaseURL)
//     if err != nil {
//         return fmt.Errorf("failed to create migrate instance: %w", err)
//     }
//     defer m.Close()

//     if err := m.Up(); err != nil && err != migrate.ErrNoChange {
//         return fmt.Errorf("failed to run migrations: %w", err)
//     }

//     return nil
// }

// for docker

// db/run.go
// func RunMigrations(databaseURL string) error {
// 	// Add retry logic for database connection
// 	maxRetries := 5
// 	var err error

// 	for i := 0; i < maxRetries; i++ {
// 		// Try to connect to database
// 		db, err := sql.Open("postgres", databaseURL)
// 		if err == nil {
// 			err = db.Ping()
// 			if err == nil {
// 				db.Close()
// 				break
// 			}
// 		}

// 		fmt.Printf("Attempt %d: Failed to connect to database, retrying in 5 seconds...\n", i+1)
// 		time.Sleep(5 * time.Second)
// 	}

// 	// Find migrations directory
// 	currentDir, err := os.Getwd()
// 	if err != nil {
// 		return fmt.Errorf("failed to get current directory: %w", err)
// 	}

// 	// In Docker environment, migrations should be in /app/db/migrations
// 	migrationsPath := "/app/db/migrations"
// 	if os.Getenv("APP_ENV") != "docker" {
// 		migrationsPath = filepath.Join(currentDir, "db", "migrations")
// 	}

// 	// Verify migrations directory exists
// 	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
// 		return fmt.Errorf("migrations directory not found at %s: %w", migrationsPath, err)
// 	}

// 	// Convert to file URL format
// 	migrationURL := fmt.Sprintf("file://%s", migrationsPath)

// 	m, err := migrate.New(migrationURL, databaseURL)
// 	if err != nil {
// 		return fmt.Errorf("failed to create migrate instance: %w", err)
// 	}
// 	defer m.Close()

// 	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 		return fmt.Errorf("failed to run migrations: %w", err)
// 	}

// 	return nil
// }

// Connect establishes a connection to the database using pgxpool
func Connect(databaseURL string) (*pgxpool.Pool, error) {

	fmt.Print("golang/ecomm-grpc/db/run.go Connect databaseURL", databaseURL)
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

// RunMigrations handles database migrations for both Docker and local environments
func RunMigrations(databaseURL string) error {
	// Add retry logic for database connection
	maxRetries := 5
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err := sql.Open("postgres", databaseURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				db.Close()
				break
			}
		}

		fmt.Printf("Attempt %d: Failed to connect to database, retrying in 5 seconds...\n", i+1)
		time.Sleep(5 * time.Second)
	}

	// if err != nil {
	// 	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	// }

	// Determine migrations path based on environment
	var migrationsPath string
	if os.Getenv("APP_ENV") == "docker" {
		migrationsPath = "/app/db/migrations"
	} else {
		// For local development
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Try different possible locations for migrations
		possiblePaths := []string{
			filepath.Join(currentDir, "db", "migrations"),
			filepath.Join(currentDir, "ecomm-grpc", "db", "migrations"),
			filepath.Join(filepath.Dir(currentDir), "ecomm-grpc", "db", "migrations"),
		}

		migrationsFound := false
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				migrationsPath = path
				migrationsFound = true
				break
			}
		}

		if !migrationsFound {
			return fmt.Errorf("migrations directory not found in any of the expected locations")
		}
	}

	// Convert to file URL format
	migrationURL := fmt.Sprintf("file://%s", migrationsPath)

	// Create and run migrations
	m, err := migrate.New(migrationURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}