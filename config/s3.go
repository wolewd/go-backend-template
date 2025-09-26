package config

import (
	"fmt"
	"log"

	"go-template/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectS3() (*minio.Client, error) {
	endpoint := utils.GetEnv("S3_ENDPOINT", "")
	accessKey := utils.GetEnv("S3_ACCESS_KEY", "")
	secretKey := utils.GetEnv("S3_SECRET_KEY", "")
	useSSL := utils.GetEnv("S3_SSL", "false")

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("S3_ENDPOINT, S3_ACCESS_KEY, and S3_SECRET_KEY must be set in environment variables")
	}

	ssl := false
	if useSSL == "true" {
		ssl = true
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: ssl,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	log.Println("S3 client initialized successfully")
	return minioClient, nil
}
