package main

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go"
)

const (
	timeout           = time.Second * 5
	spacesURLTemplate = "https://%s.%s/%s"
)

// Storage is used to save/retrive dump file from remote object storage
type Storage interface {
	Save(ctx context.Context, file *TempFile) (string, error)
}

type SpacesConfig struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
}

type SpacesStorage struct {
	client *minio.Client
	cfg    SpacesConfig
}

func NewSpacesStorage(cfg SpacesConfig) (*SpacesStorage, error) {
	client, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, false)
	if err != nil {
		return nil, err
	}

	return &SpacesStorage{
		client: client,
		cfg:    cfg,
	}, nil
}

func (s *SpacesStorage) Save(ctx context.Context, file *TempFile) (string, error) {
	opts := minio.PutObjectOptions{
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	size, err := file.Size()
	if err != nil {
		return "", err
	}
	_, err = s.client.PutObjectWithContext(ctx,
		s.cfg.Bucket, file.Name(), file.Reader(), size, opts)
	if err != nil {
		return "", err
	}

	return s.generateFileURL(file.Name()), nil
}

// DigitalOcean Spaces link format
func (s *SpacesStorage) generateFileURL(filename string) string {
	return fmt.Sprintf(spacesURLTemplate, s.cfg.Bucket, s.cfg.Endpoint, filename)
}
