package spacer

import (
	"context"
	"fmt"
	"github.com/minio/minio-go"
)

const (
	// DigitalOcean Spaces link format
	spacesURLTemplate = "https://%s.%s/%s"
)

type SpacesStorage struct {
	client   *minio.Client
	endpoint string
	bucket   string
}

func NewSpacesStorage(endpoint, bucket, accessKey, secretKey string) (*SpacesStorage, error) {
	client, err := minio.New(endpoint, accessKey, secretKey, false)
	if err != nil {
		return nil, err
	}

	return &SpacesStorage{
		client:   client,
		endpoint: endpoint,
		bucket:   bucket,
	}, nil
}

func (s *SpacesStorage) Save(ctx context.Context, file *TempFile) (string, error) {
	opts := minio.PutObjectOptions{
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	size, err := file.Size()
	if err != nil {
		return "", err
	}
	_, err = s.client.PutObjectWithContext(ctx,
		s.bucket, file.Name(), file.Reader(), size, opts)
	if err != nil {
		return "", err
	}

	return s.generateFileURL(file.Name()), nil
}

func (s *SpacesStorage) GetLatest(ctx context.Context) (string, error) {
	return "", nil
}

func (s *SpacesStorage) generateFileURL(filename string) string {
	return fmt.Sprintf(spacesURLTemplate, s.bucket, s.endpoint, filename)
}
