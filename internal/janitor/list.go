package janitor

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Object struct {
	Bucket string
	Size   int64
	Key    string
}

func (j *Janitor) listObjects() {
	defer j.wg.Done()
	defer close(j.objectChan)

	var wg sync.WaitGroup

	for _, bucket := range j.buckets {
		wg.Add(1)
		j.objectSem <- token{}

		go j.listBucket(bucket, &wg)
	}

	wg.Wait()
}

const batchSize int = 1000

func (j *Janitor) listBucket(bucket string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-j.objectSem }()

	ctx := context.Background()

	responseLength := batchSize
	afterKey := ""
	batchSizeInt32 := int32(batchSize)

	for responseLength == batchSize {
		response, err := j.client.ListObjectsV2(ctx,
			&s3.ListObjectsV2Input{
				Bucket:     aws.String(bucket),
				StartAfter: &afterKey,
				MaxKeys:    batchSizeInt32,
			})
		if err != nil {
			log.Printf("failed listing bucket %s: %s\n", bucket, err)
			return
		}

		responseLength = len(response.Contents)
		if responseLength == batchSize {
			afterKey = *response.Contents[responseLength-1].Key
		}

		for _, object := range response.Contents {
			j.objectChan <- Object{
				Bucket: bucket,
				Size:   object.Size,
				Key:    *object.Key,
			}
		}
	}
}
