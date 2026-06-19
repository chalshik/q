package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	migrationsfs "dispatcher/db"
)

// NewPostgresConnection initializes the database pool using environment variables or defaults.
func NewPostgresConnection() (*sqlx.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASSWORD", "secret")
	name := getEnv("DB_NAME", "dispatcher")
	ssl := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, name, ssl)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	// Set connection pool baselines
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection works
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	// Run migrations automatically on successful connection
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("migration runner failed: %w", err)
	}

	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	// Wrap the embedded filesystem from our migrations package
	sourceDriver, err := iofs.New(migrationsfs.FS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration iofs source: %w", err)
	}

	dbDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to construct migrator: %w", err)
	}

	// Execute up migrations
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("migrations: database already up to date, no changes applied")
			return nil
		}
		return fmt.Errorf("failed to execute migrations up: %w", err)
	}

	log.Println("migrations: applied successfully")
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
