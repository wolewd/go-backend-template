package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type Config struct {
	DB *pgxpool.Pool
	S3 *minio.Client
}

func InitConfig() (*Config, error) {
	// init db
	database, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	// init s3
	s3Client, err := ConnectS3()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB: database,
		S3: s3Client,
	}, nil
}
