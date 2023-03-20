package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Port         int                `json:"port"`
	CloudStorage CloudStorageConfig `json:"cloud_storage"`
	Firestore    FirestoreConfig    `json:"firestore"`
	Postgres     PostgresConfig     `json:"postgres"`
}

type CloudStorageConfig struct {
	BucketName string `json:"bucket_name"`
}

type FirestoreConfig struct {
	ProjectID      string `json:"project_id"`
	CollectionName string `json:"collection_name"`
}

type PostgresConfig struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"database_name"`
	TableName    string `json:"table_name"`
}

func New(configFile string) (*Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, err

}
