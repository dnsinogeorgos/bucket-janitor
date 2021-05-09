package config

import (
	"encoding/json"
	"os"
)

// Config stores the configuration loaded during startup
type Config struct {
	AwsAccessKeyId string   `json:"aws_access_key_id"`
	AwsSecretKey   string   `json:"aws_secret_key"`
	S3Region       string   `json:"s3_region"`
	S3Buckets      []string `json:"s3_buckets"`
}

// NewConfig creates a config instance and loads json content from file
func NewConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return &Config{}, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return &Config{}, err
	}

	return config, nil
}
