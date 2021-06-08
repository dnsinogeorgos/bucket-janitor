package aws

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	ac "github.com/aws/aws-sdk-go-v2/config"
)

// Config stores the configuration loaded during startup
type Config struct {
	Client      *s3.Client
	Downloader  *manager.Downloader
	S3Buckets   []string `json:"s3_buckets,omitempty"`
	Concurrency int      `json:"concurrency,omitempty"`
}

// CreateLoad creates a config instance, loads json content from file and initializes
// AWS client and downloader
func CreateLoad(cp string) *Config {
	c := &Config{}

	f, err := os.Open(cp)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := ac.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	c.Client = s3.NewFromConfig(cfg)
	c.Downloader = manager.NewDownloader(c.Client)

	return c
}
