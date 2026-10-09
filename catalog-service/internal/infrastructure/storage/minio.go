package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Lama189/soundwave-platform/catalog-service/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client       *minio.Client
	bucketName   string
	uploadExpiry time.Duration
}

func NewMinioStorage(ctx context.Context, cfg config.StorageConfig) (*MinioStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio client: %w", err)
	}

	storage := &MinioStorage{
		client:       client,
		bucketName:   cfg.BucketName,
		uploadExpiry: cfg.UploadExpiry,
	}

	if err := storage.ensureBucket(ctx); err != nil {
		return nil, fmt.Errorf("ensure bucket exists: %w", err)
	}

	return storage, nil
}

func (s *MinioStorage) GenerateUploadUrl(ctx context.Context, key string) (string, error) {
	presignedURL, err := s.client.PresignedPutObject(ctx, s.bucketName, key, s.uploadExpiry)
	if err != nil {
		return "", fmt.Errorf("generate presigned put url for %s: %w", key, err)
	}

	return presignedURL.String(), nil
}

func (s *MinioStorage) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("check bucket %s: %w", s.bucketName, err)
	}

	if !exists {
		err := s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("create bucket %s: %w", s.bucketName, err)
		}
	}

	return nil
}
