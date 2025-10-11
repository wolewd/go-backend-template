package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {
	dbuser := GetEnv("DB_USER", "")
	dbpassword := GetEnv("DB_PASSWORD", "")
	dbhost := GetEnv("DB_HOST", "")
	dbport := GetEnv("DB_PORT", "")
	dbname := GetEnv("DB_NAME", "")
	dbsslmode := GetEnv("DB_SSL_MODE", "disable")

	if dbuser == "" || dbpassword == "" || dbhost == "" || dbport == "" || dbname == "" {
		return nil, fmt.Errorf("DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, and DB_NAME must be set in environment variables")
	}

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbuser, dbpassword, dbhost, dbport, dbname, dbsslmode,
	)

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connection pool failed: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	log.Println("Database connected successfully")
	return pool, nil
}
