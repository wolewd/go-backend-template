package utils

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

type S3Helper struct {
	client *minio.Client
	ctx    context.Context
}

func NewS3Helper(client *minio.Client) *S3Helper {
	return &S3Helper{
		client: client,
		ctx:    context.Background(),
	}
}

func (s *S3Helper) UploadFile(bucketName, objectPath string, data []byte, contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(s.ctx, bucketName, objectPath, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return objectPath, nil
}

func (s *S3Helper) UploadFileWithUUID(bucketName, originalFileName string, data []byte, contentType string) (string, error) {
	ext := filepath.Ext(originalFileName)
	newFileName := GenerateUUIDv7() + ext

	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(s.ctx, bucketName, newFileName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return newFileName, nil
}

func (s *S3Helper) FileExists(bucketName, objectPath string) (bool, error) {
	_, err := s.client.StatObject(s.ctx, bucketName, objectPath, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat object: %w", err)
	}
	return true, nil
}

func (s *S3Helper) DeleteFile(bucketName, objectPath string) error {
	return s.client.RemoveObject(s.ctx, bucketName, objectPath, minio.RemoveObjectOptions{})
}

func (s *S3Helper) DownloadPublicFile(bucketName, objectPath string) string {
	endpoint := GetEnv("S3_ENDPOINT", "")
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectPath)
}

func (s *S3Helper) DownloadPrivateFile(bucketName, objectPath string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(s.ctx, bucketName, objectPath, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return presignedURL.String(), nil
}
