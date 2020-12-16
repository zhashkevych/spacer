package spacer

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go"
	"io/ioutil"
	"net/http"
)

const (
	// DigitalOcean Spaces link format
	spacesURLTemplate = "https://%s.%s/%s"
)

// SpacesStorage is a DigitalOcean Spaces client
type SpacesStorage struct {
	client   *minio.Client
	endpoint string
	bucket   string
}

// NewSpacesStorage creates new DigitalOcean Spaces client
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

// Save saves files to Digital Ocean Spaces
func (s *SpacesStorage) Save(ctx context.Context, file *DumpFile) (string, error) {
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

// GetLatest downloads
func (s *SpacesStorage) GetLatest(ctx context.Context, prefix string) (*DumpFile, error) {
	name, err := s.getLatestDumpName(ctx, prefix)
	if err != nil {
		return nil, err
	}

	url := s.generateFileURL(name)
	fileData, err := s.fetch(url)
	if err != nil {
		return nil, err
	}

	return s.createTempFile(fileData)
}

func (s *SpacesStorage) getLatestDumpName(ctx context.Context, prefix string) (string, error) {
	objects, err := s.parseObjects(prefix)
	if err != nil {
		return "", err
	}

	latestObject := s.getLatestObject(objects)

	return latestObject.Key, nil
}

func (s *SpacesStorage) parseObjects(prefix string) ([]minio.ObjectInfo, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	objects := make([]minio.ObjectInfo, 0)
	objectsCh := s.client.ListObjects(s.bucket, prefix, true, doneCh)
	for object := range objectsCh {
		if object.Err != nil {
			continue
		}

		objects = append(objects, object)
	}

	if len(objects) == 0 {
		return nil, errors.New("no files found")
	}

	return objects, nil
}

func (s *SpacesStorage) getLatestObject(objects []minio.ObjectInfo) minio.ObjectInfo {
	latestObject := objects[0]
	for _, object := range objects {
		if latestObject.LastModified.Unix() < object.LastModified.Unix() {
			latestObject = object
		}
	}
	return latestObject
}

func (s *SpacesStorage) fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (* SpacesStorage) createTempFile(data []byte) (*DumpFile, error) {
	tempFile, err := NewDumpFile("restore")
	if err != nil {
		return nil, err
	}

	if err := tempFile.Write(data); err != nil {
		return nil, err
	}

	return tempFile, nil
}

func (s *SpacesStorage) generateFileURL(filename string) string {
	return fmt.Sprintf(spacesURLTemplate, s.bucket, s.endpoint, filename)
}
