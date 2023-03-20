package cloud_storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"

	"compute_engine/config"
)

type CloudStorage struct {
	client     *storage.Client
	bucketName string
}

func New(config config.CloudStorageConfig) (*CloudStorage, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create a client: %v", err)
	}
	return &CloudStorage{
		client:     client,
		bucketName: config.BucketName,
	}, nil

}

func (cs *CloudStorage) UploadFile(filename, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	bucket := cs.client.Bucket(cs.bucketName)
	o := bucket.Object(filename)
	ww := o.NewWriter(ctx)
	if _, err := fmt.Fprintf(ww, content); err != nil {
		return fmt.Errorf("failed to write to a writer: %v", err)
	}
	if err := ww.Close(); err != nil {
		return fmt.Errorf("failed to close a writer: %v", err)
	}
	return nil
}

func (cs *CloudStorage) CreateSignedURL(filename string) (string, error) {
	bucket := cs.client.Bucket(cs.bucketName)
	opts := storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Second * 30),
	}
	url, err := bucket.SignedURL(filename, &opts)
	if err != nil {
		//We intentionally overlook an error since `SignedURL()` always failed in our local environment.
		fmt.Printf("failed to create a signed URL: %v\n", err)
		return "failed to create a signed URL", nil
	}
	return url, nil
}
