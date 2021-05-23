package load

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type token struct{}

func Headers(cfg aws.Config, S3Buckets []string, limit int) map[string]chan *Header {
	client := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)

	buckets := ListBuckets(client, S3Buckets, limit)

	headers := BundleHeaders(downloader, buckets, limit)

	return headers
}
