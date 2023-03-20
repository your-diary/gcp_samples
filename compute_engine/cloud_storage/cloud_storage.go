package cloud_storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"

	"compute_engine/config"
)

func UploadFile(config config.CloudStorageConfig, filename, content string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create a client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	bucket := client.Bucket(config.BucketName)
	o := bucket.Object(filename)
	ww := o.NewWriter(ctx)
	if _, err = fmt.Fprintf(ww, content); err != nil {
		return "", fmt.Errorf("failed to write to a writer: %v", err)
	}
	if err := ww.Close(); err != nil {
		return "", fmt.Errorf("failed to close a writer: %v", err)
	}

    //TODO
	// 	opts := storage.SignedURLOptions{
	// 		Method:  "GET",
	// 		Expires: time.Now().Add(time.Second * 30),
	// 	}
	// 	url, err := bucket.SignedURL(filename, &opts)
	// 	if err != nil {
	// 		return "", fmt.Errorf("failed to create a signed URL: %v", err)
	// 	}
	url := "url"
	return url, nil
}
