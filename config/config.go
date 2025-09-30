package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DB *pgxpool.Pool
}

func InitConfig() (*Config, error) {
	// init db
	database, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB: database,
	}, nil
}
