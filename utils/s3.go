package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	s3Client     *minio.Client
	s3ClientOnce sync.Once
	s3InitErr    error
)

func getS3Client() (*minio.Client, error) {
	s3ClientOnce.Do(func() {
		endpoint := GetEnv("S3_ENDPOINT", "")
		accessKey := GetEnv("S3_ACCESS_KEY", "")
		secretKey := GetEnv("S3_SECRET_KEY", "")
		useSSL := GetEnv("S3_SSL", "false")

		if endpoint == "" || accessKey == "" || secretKey == "" {
			s3InitErr = fmt.Errorf("S3_ENDPOINT, S3_ACCESS_KEY, and S3_SECRET_KEY must be set in environment variables")
			return
		}

		ssl := false
		if useSSL == "true" {
			ssl = true
		}

		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: ssl,
		})

		if err != nil {
			s3InitErr = fmt.Errorf("failed to initialize S3 client: %w", err)
			return
		}

		s3Client = client
		log.Println("S3 client initialized successfully")
	})

	if s3InitErr != nil {
		return nil, s3InitErr
	}

	return s3Client, nil
}

func buildPath(prefix, fileName string) string {
	if prefix == "" {
		return fileName
	}
	return filepath.Join(prefix, fileName)
}

// Example: UploadFile("bucket", "uploads/images", "photo.jpg", data, "image/jpeg")
// Result: uploads/images/photo.jpg
// prefix is optional, use "" for root level
func UploadFile(bucketName, prefix, fileName string, data []byte, contentType string) (string, error) {
	client, err := getS3Client()
	if err != nil {
		return "", err
	}

	objectPath := buildPath(prefix, fileName)

	reader := bytes.NewReader(data)
	_, err = client.PutObject(context.Background(), bucketName, objectPath, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return objectPath, nil
}

// Example: UploadFileWithUUID("bucket", "users/avatars", "photo.jpg", data, "image/jpeg")
// Result: users/avatars/01936b3e-4d2a-7890-abcd-ef1234567890.jpg
// prefix is optional, use "" for root level
func UploadFileWithUUID(bucketName, prefix, originalFileName string, data []byte, contentType string) (string, error) {
	client, err := getS3Client()
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(originalFileName)
	uuidFileName := GenerateUUIDv7() + ext
	objectPath := buildPath(prefix, uuidFileName)

	reader := bytes.NewReader(data)
	_, err = client.PutObject(context.Background(), bucketName, objectPath, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return objectPath, nil
}

// Example: FileExists("bucket", "uploads/images", "photo.jpg")
// prefix is optional, use "" for root level
func FileExists(bucketName, prefix, fileName string) (bool, error) {
	client, err := getS3Client()
	if err != nil {
		return false, err
	}

	objectPath := buildPath(prefix, fileName)

	_, err = client.StatObject(context.Background(), bucketName, objectPath, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat object: %w", err)
	}
	return true, nil
}

// Example: DeleteFile("bucket", "uploads/images", "photo.jpg")
// prefix is optional, use "" for root level
func DeleteFile(bucketName, prefix, fileName string) error {
	client, err := getS3Client()
	if err != nil {
		return err
	}

	objectPath := buildPath(prefix, fileName)
	return client.RemoveObject(context.Background(), bucketName, objectPath, minio.RemoveObjectOptions{})
}

// Example: DownloadPublicFile("bucket", "uploads/images", "photo.jpg")
// prefix is optional, use "" for root level
func DownloadPublicFile(bucketName, prefix, fileName string) string {
	endpoint := GetEnv("S3_ENDPOINT", "")
	useSSL := GetEnv("S3_SSL", "false")

	protocol := "http"
	if useSSL == "true" {
		protocol = "https"
	}

	objectPath := buildPath(prefix, fileName)
	return fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, bucketName, objectPath)
}

// Example: DownloadPrivateFile("bucket", "private/docs", "document.pdf", time.Hour)
// prefix is optional, use "" for root level
func DownloadPrivateFile(bucketName, prefix, fileName string, expiry time.Duration) (string, error) {
	client, err := getS3Client()
	if err != nil {
		return "", err
	}

	objectPath := buildPath(prefix, fileName)
	reqParams := make(url.Values)
	presignedURL, err := client.PresignedGetObject(context.Background(), bucketName, objectPath, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return presignedURL.String(), nil
}

// GetS3Client returns the initialized S3 client for advanced usage
func GetS3Client() (*minio.Client, error) {
	return getS3Client()
}
