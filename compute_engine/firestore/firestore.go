package firestore

import (
	"compute_engine/config"
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Firestore struct {
	client         *firestore.Client
	collectionName string
}

func New(config config.FirestoreConfig) (*Firestore, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create a client: %v", err)
	}
	return &Firestore{
		client:         client,
		collectionName: config.CollectionName,
	}, nil
}

func (fs *Firestore) Insert(content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, _, err := fs.client.Collection(fs.collectionName).Add(ctx, map[string]interface{}{
		"timestamp": time.Now(),
		"content":   content,
	})
	return err
}

func (fs *Firestore) selectByContent(content string) ([][]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	iter := fs.client.Collection(fs.collectionName).Where("content", "==", content).Documents(ctx)
	var ret [][]any
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, []any{
			doc.Data()["timestamp"],
			doc.Data()["content"],
		})
	}
	return ret, nil
}
