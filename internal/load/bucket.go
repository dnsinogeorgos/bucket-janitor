package load

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ListBuckets returns a chan of s3 object slices
func ListBuckets(client *s3.Client, S3Buckets []string, limit int) map[string]<-chan types.Object {
	sem := make(chan token, limit)
	channels := make(map[string]<-chan types.Object)

	for _, bucket := range S3Buckets {
		channels[bucket] = ListBucket(client, bucket, sem)
	}

	return channels
}

// batchSize is the number of objects to list on each API call to S3
const batchSize int = 1000

// ListBucket returns a chan of s3 objects and ships a goroutine to fetch them
func ListBucket(client *s3.Client, bucket string, sem chan token) <-chan types.Object {
	c := make(chan types.Object)

	go func() {
		sem <- token{}
		go func() {
			ctx := context.Background()

			responseLength := batchSize
			afterKey := ""
			batchSizeInt32 := int32(batchSize)

			for responseLength == batchSize {
				response, err := client.ListObjectsV2(ctx,
					&s3.ListObjectsV2Input{
						Bucket:     aws.String(bucket),
						StartAfter: &afterKey,
						MaxKeys:    batchSizeInt32,
					})
				if err != nil {
					log.Printf("failed listing bucket %s: %s\n", bucket, err)
					close(c)
					<-sem
					return
				}

				responseLength = len(response.Contents)
				if responseLength == batchSize {
					afterKey = *response.Contents[responseLength-1].Key
				}

				for _, object := range response.Contents {
					c <- object
				}
			}

			close(c)
			<-sem
		}()
	}()

	return c
}
