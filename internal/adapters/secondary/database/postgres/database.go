package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func ConnectDB() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	print("Connecting to database with the following parameters:\n")
	print("Host: ", host, "\n")
	print("Port: ", port, "\n")
	print("User: ", user, "\n")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("unable to parse DSN: %w", err)
	}

	// Ajustes para a Rinha (consumo controlado)
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 15 * time.Minute
	config.HealthCheckPeriod = 5 * time.Minute

	pool, err := pgxpool.New(context.Background(), config.ConnString())
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	dbPool = pool
	return nil
}

func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
	}
}

func GetDB() *pgxpool.Pool {
	return dbPool
}
