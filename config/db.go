package config

import (
	"context"
	"fmt"
	"log"

	"go-template/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {
	dbuser := utils.GetEnv("DB_USER", "")
	dbpassword := utils.GetEnv("DB_PASSWORD", "")
	dbhost := utils.GetEnv("DB_HOST", "")
	dbport := utils.GetEnv("DB_PORT", "")
	dbname := utils.GetEnv("DB_NAME", "")
	dbsslmode := utils.GetEnv("DB_SSL_MODE", "disable")

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
